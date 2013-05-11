package main

import (
	"flag"
	"net/http"

	"github.com/dustin/go-couchbase"
)

var bucket *couchbase.Bucket

func main() {
	addr := flag.String("addr", ":8482", "")

	flag.Parse()

	var err error
	bucket, err = couchbase.GetBucket("http://localhost:8091/", "default", "l4d2")
	if err != nil {
		panic(err)
	}
	defer bucket.Close()

	panic(http.ListenAndServe(*addr, nil))
}
