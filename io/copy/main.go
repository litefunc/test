package main

import (
	"VodoPlay/logger"
	"encoding/json"
	"io"
	"os"
	"strings"
)

func main() {
	var is []int
	for i := 0; i < 10; i++ {
		is = append(is, i)
	}
	by, _ := json.Marshal(is)
	// ioutil.WriteFile("file", by, os.ModePerm)

	src := strings.NewReader(string(by))

	f, err := os.OpenFile("file", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return
	}
	defer f.Close()

	// var w writer
	if _, err := io.Copy(f, src); err != nil {
		logger.Error(err)
	}

	copyFile("file", "file1")
}

type writer struct{}

func (rec writer) Write(p []byte) (n int, err error) {
	logger.Debug(len(p), string(p))
	return len(p), nil
}

func copyFile(oldpath, newpath string) error {

	f, err := os.Stat(oldpath)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Debug(f.ModTime())

	f1, err := os.Open(oldpath)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer f1.Close()

	// f2, err := os.Create(newpath)

	if err := os.RemoveAll(newpath); err != nil {
		logger.Error(err)
		return err
	}

	f2, err := os.OpenFile(newpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
	if err != nil {
		logger.Error(err)
		return err
	}
	defer f2.Close()

	if _, err := io.Copy(f2, f1); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
