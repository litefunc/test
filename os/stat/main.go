package main

import (
	"LocalServer/lib/logger"
	"os"
)

func main() {
	p := "/home/david/program/go/src/test/os/statdir/dir/test.json"
	// p := "/home/david/program/go/src/LocalServer/test/testdata/tmp/VOD-01EWJ6JPP8Q1AKAPMN862YZ724/cover/cover.json"
	_, err := os.Stat(p)
	if err != nil && !os.IsNotExist(err) {
		logger.Error(err)
	}
}
