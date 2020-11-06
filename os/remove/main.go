package main

import (
	"VodoPlay/logger"
	"os"
)

func main() {
	if err := os.RemoveAll("dir"); err != nil {
		logger.Error(err)
	}
	if err := os.RemoveAll("file"); err != nil {
		logger.Error(err)
	}
	if err := os.Remove("file"); err != nil {
		logger.Error(err)
	}
	if err := os.Remove("file1"); err != nil {
		logger.Error(err)
	}
}
