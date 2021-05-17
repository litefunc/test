package test

import (
	"path"
	"runtime"
)

var rootDir string

func RootDir() string {
	return rootDir
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir = path.Dir(filename)
}
