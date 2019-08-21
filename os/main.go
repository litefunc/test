package main

import (
	"cloud/lib/logger"
	"os"
	"path/filepath"
)

func main() {
	path := os.Getenv("GOPATH") + "/src/test/os/testdir/docker.json"

	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		logger.Error(err)
		return
	}

	// if err := os.RemoveAll(dir); err != nil {
	// 	logger.Error(err)
	// }
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logger.Error(err)
	}
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

}
