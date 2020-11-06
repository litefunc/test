package main

import (
	"VodoPlay/logger"
	"os"
)

func main() {
	// rename("file1", "dir1/file1")
	rename("dir2", "dir1/dir")
}

func rename(oldpath, newpath string) error {
	if err := os.Rename(oldpath, newpath); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
