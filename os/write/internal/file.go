package internal

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"test/logger"
	"time"
)

func openFiles(path string, n int) []*os.File {
	fs := make([]*os.File, n, n)
	for i := 0; i < n; i++ {
		f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			logger.Fatal((err))
		}
		fs[i] = f
	}
	return fs
}

func write(f *os.File, id int, n int) {

	for i := 0; i < n; i++ {
		if _, err := f.Write([]byte(strconv.Itoa(id))); err != nil {
			logger.Fatal(err)
		}
	}
}

func empty(filename string) {

	if err := ioutil.WriteFile(filename, []byte{}, os.ModePerm); err != nil {
		logger.Fatal(err)
	}
}

func write1(filename string, n int) {

	s := strings.Repeat("1", n)
	if err := ioutil.WriteFile(filename, []byte(s), os.ModePerm); err != nil {
		logger.Fatal(err)
	}
}

// check if file lock automatically when write
func CoWrite() {
	p := "test.txt"
	fs := openFiles(p, 10)
	var wg sync.WaitGroup
	for i := range fs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Millisecond)
			write(fs[i], i, 1000)
		}(i)
	}
	wg.Wait()

	for i := range fs {
		fs[i].Close()
	}

	fi, err := os.Stat(p)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug(fi.Size())
	check(p)
}

// check if file lock automatically when read
func CoReadWrite() {
	p := "test.txt"
	if err := os.Truncate(p, 0); err != nil {
		logger.Fatal(err)
	}

	n := 10000
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			write1(p, n)
			empty(p)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			check1(p, n)
		}
	}()

	wg.Wait()

	fi, err := os.Stat(p)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug(fi.Size())
	check(p)
}

func check(path string) {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal(err)
	}

	for i := range by {
		if by[i] != by[0] {
			logger.Fatal(string(by[i]))
		}
	}
}

func check1(path string, n int) {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal(err)
	}
	ln := len(by)

	if ln != 0 && ln != n {
		logger.Fatal(ln)
	}
}
