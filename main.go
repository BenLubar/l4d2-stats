package main

import (
	"os"
	"fmt"
	"io"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	read := LogReader(f)

	for {
		ln, line, err := read()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("ERROR:%d: %s\n", ln, err)
			}
			return
		}

		fmt.Printf("%d %#v\n", ln, line)
	}
}
