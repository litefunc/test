package internal

import (
	"io/ioutil"
	"os"
	"test/logger"
)

func Temp() {
	s, err := ioutil.TempDir("dir", "abc")
	logger.LogErr(err)

	dir, err := os.Getwd()
	s, err = ioutil.TempDir(dir, "abc")
	logger.LogErr(err)
	logger.Debug(s)
	defer os.RemoveAll(s)

	f, err := ioutil.TempFile(s, "")
	logger.LogErr(err)
	logger.Debug(f.Name())

	f, err = ioutil.TempFile(s, "foo*.txt")
	logger.LogErr(err)
	logger.Debug(f.Name())

}
