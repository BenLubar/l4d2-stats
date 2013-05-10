package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dustin/go-couchbase"
)

var (
	offsets    = make(map[string]int)
	offsetLock sync.Mutex

	bucket *couchbase.Bucket
)

func processFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	base := filepath.Base(filename)

	read := LogReader(f)

	offsetLock.Lock()
	off := offsets[base]
	offsetLock.Unlock()

	for i := 0; i < off; i++ {
		_, _, err := read()
		if err != nil {
			return err
		}
	}

	for {
		ln, line, err := read()
		if err != nil {
			if err == io.EOF {
				offsetLock.Lock()
				offsets[base] = ln - 1
				offsetLock.Unlock()
				return nil
			}
			return err
		}

		line["Base"] = base
		_, err = bucket.Add(fmt.Sprintf("%s+%d", base, ln), 0, line)
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
	workers := flag.Int("workers", 4, "")
	watch := flag.Bool("watch", false, "")

	flag.Parse()

	var err error
	bucket, err = couchbase.GetBucket("http://localhost:8091/", "default", "default")
	if err != nil {
		panic(err)
	}
	defer bucket.Close()

	wg := &sync.WaitGroup{}

	ch := make(chan string)

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go processFiles(ch, wg)
	}

	if *watch {
		watcher(ch, flag.Args())
	} else {
		for _, path := range flag.Args() {
			abspath, err := filepath.Abs(path)
			if err != nil {
				log.Fatalf("Error absing %v: %v", path, err)
			}
			ch <- abspath
		}
		close(ch)

		wg.Wait()
	}
}
