package main

import (
	"cloud/lib/logger"
	"os"
	"path/filepath"
	"strings"
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

	stat(path, "/usr/", "usr/abc", "Behind%20My%20Life/Behind%20My%20Life.jpg")
	p := os.Getenv("GOPATH") + "/src/test/os/" + "Behind%20My%20Life/Behind%20My%20Life.jpg"
	stat(p)
	stat(strings.Replace(p, "%20", " ", -1))
}

func stat(path ...string) {
	for _, p := range path {
		if _, err := os.Stat(p); err != nil {
			logger.Error(err)
		}
	}
}
