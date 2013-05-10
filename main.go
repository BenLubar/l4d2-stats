package main

import (
	"fmt"
	"github.com/dustin/go-couchbase"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	bucket, err := couchbase.GetBucket("http://localhost:8091/", "default", "default")
	if err != nil {
		panic(err)
	}
	defer bucket.Close()

	var wg sync.WaitGroup
	wg.Add(len(os.Args) - 1)

	for i := 1; i < len(os.Args); i++ {
		go func(i int) {
			defer wg.Done()
			filename, err := filepath.Abs(os.Args[i])
			if err != nil {
				panic(err)
			}

			f, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			base := filepath.Base(filename)

			read := LogReader(f)

			for {
				ln, line, err := read()
				if err != nil {
					if err != io.EOF {
						panic(err)
					}
					return
				}

				line["Base"] = base
				_, err = bucket.Add(fmt.Sprintf("%s+%d", base, ln), 0, line)
				if err != nil {
					panic(err)
				}
			}
		}(i)
	}

	wg.Wait()
}
