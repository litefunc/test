package main

import (
	"cloud/lib/logger"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"test/http/client/cmd/range/internal"
)

type header struct {
	AcceptRanges  bool
	ContentLength int
}

func getHeader(url string) (header, error) {

	var h header

	res, err := http.Head(url)
	if err != nil {
		logger.Error(err)
		return h, err
	}

	i, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		logger.Error(err)
		return h, err
	}
	h.ContentLength = i

	if res.Header.Get("Accept-Ranges") == "bytes" {
		h.AcceptRanges = true
	}

	return h, nil

}

func rangeGet(url string, start, end int) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// logger.Debug(string(by))

	return by, nil
}

func get(url string) ([]byte, error) {
	t1 := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// logger.Debug(string(by))

	diff := time.Now().Sub(t1)
	logger.Debug(diff)

	return by, nil
}

func RangeGet(url string, total, size int) ([]byte, error) {

	var all []byte

	n := total / size

	for i := 0; i < n; i++ {
		by, err := rangeGet(url, i*size, i*size+size-1)
		if err != nil {
			return nil, err
		}
		all = append(all, by...)
	}
	r := total % size
	if r > 0 {
		by, err := rangeGet(url, n*size, total-1)
		if err != nil {
			return nil, err
		}
		all = append(all, by...)
	}

	return all, nil
}

func Get(url string) ([]byte, error) {

	t1 := time.Now()
	h, err := getHeader(url)
	if err != nil {
		return nil, err
	}
	if !h.AcceptRanges {
		return get(url)
	}

	size := 1048576 * 100
	if h.ContentLength <= size {
		return get(url)
	}

	by, err := RangeGet(url, h.ContentLength, size)
	if err != nil {
		return nil, err
	}
	diff := time.Now().Sub(t1)
	logger.Debug(diff)
	return by, nil
}

func main() {

	// cli := internal.NewClient(1048576 * 100)

	us := []string{}
	var b []byte
	var err error

	b, err = get(us[1])
	logger.Debug(len(b), err)
	// logger.Debug(string(b))

	// b, err = cli.Get(us[1])
	// logger.Debug(len(b), err)
	// logger.Debug(string(b))

	internal.RangeGet(us[1])
	logger.Debug(len(b), err)

}
