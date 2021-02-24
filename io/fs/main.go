package main

import (
	"io/fs"
	"test/logger"
)

func main() {
	logger.Debug(fs.ModeDir, uint32(fs.ModeDir))
	logger.Debug(fs.ModePerm, uint32(fs.ModePerm))
}
