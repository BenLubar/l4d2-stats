package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dustin/go-couchbase"
)

var bucket *couchbase.Bucket

func processFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	base := filepath.Base(filename)

	read := LogReader(f)

	for {
		ln, line, err := read()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}

		line["Base"] = base
		err = bucket.Set(fmt.Sprintf("%s+%d", base, ln), 0, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func processFiles(fns <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for filename := range fns {
		err := processFile(filename)
		if err != nil {
			log.Fatalf("Error on %v: %v", filename, err)
		}
	}
}

func main() {
	var err error
	bucket, err = couchbase.GetBucket("http://localhost:8091/", "default", "default")
	if err != nil {
		panic(err)
	}
	defer bucket.Close()

	wg := &sync.WaitGroup{}

	ch := make(chan string)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go processFiles(ch, wg)
	}

	for i := 1; i < len(os.Args); i++ {
		abspath, err := filepath.Abs(os.Args[i])
		if err != nil {
			log.Fatalf("Error absing %v: %v", os.Args[i], err)
		}
		ch <- abspath
	}
	close(ch)

	wg.Wait()
}
