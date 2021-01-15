package main

import (
	"os"
	"test/logger"
)

func main() {
	f, err := os.OpenFile("a.txt", os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

	logger.Debug(os.ModePerm.String())
	logger.Debug(os.FileMode(0666).String())
}
