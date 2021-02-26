package main

import (
	"io/fs"
	"test/logger"
)

func main() {
	logger.Debugf(`%b %d %v`, fs.ModeDir, fs.ModeDir, fs.ModeDir)
	logger.Debugf(`%b %d %v`, fs.ModeAppend, fs.ModeAppend, fs.ModeAppend)
	logger.Debugf(`%b %d %v`, fs.ModeExclusive, fs.ModeExclusive, fs.ModeExclusive)
	logger.Debugf(`%b %d %v`, fs.ModeTemporary, fs.ModeTemporary, fs.ModeTemporary)
	logger.Debugf(`%b %d %v`, fs.ModeSymlink, fs.ModeSymlink, fs.ModeSymlink)
	logger.Debugf(`%b %d %v`, fs.ModeDevice, fs.ModeDevice, fs.ModeDevice)
	logger.Debugf(`%b %d %v`, fs.ModeIrregular, fs.ModeIrregular, fs.ModeIrregular)
	logger.Debugf(`%b %d %v`, fs.ModeType, fs.ModeType, fs.ModeType)
	logger.Debugf(`%b %d %v`, fs.ModePerm, fs.ModePerm, fs.ModePerm)
	logger.Debugf(`%b %d %v`, fs.FileMode(0644), fs.FileMode(0644), fs.FileMode(0644))

	logger.Trace(fs.ModeDir.Type(), fs.ModePerm.Type())
}
