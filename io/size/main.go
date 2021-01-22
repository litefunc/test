package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"test/logger"
)

func main() {
	f, err := os.OpenFile("a.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

	size(f, "a.txt", "a")
	size(f, "a.txt", "a\n")
	size(f, "a.txt", "a a")
	size(f, "a.txt", "a a ")
	size(f, "a.txt", "a\ta")
}

func size(f *os.File, filename, s string) {
	ioutil.WriteFile(filename, []byte(s), os.ModePerm)

	fi, err := f.Stat()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(fmt.Sprintf(`%q`, s), fi.Size())
}
