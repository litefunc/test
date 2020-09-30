package internal

import (
	"fmt"
	"io/ioutil"
	"mstore/logger"
	"net/http"
	"strconv"
	"time"
)

type Range struct {
	Id    uint64
	Start uint64
	End   uint64
}

type Ranges []Range

type division struct {
	total   uint64
	partial uint64
	exclude uint64
	startId uint64
}

func getRanges(d division) Ranges {

	rgs := Ranges{}

	// return if parameter invalid
	if d.total == 0 || d.partial == 0 || d.total < d.exclude {
		return rgs
	}

	id := d.startId

	start := d.exclude
	end := d.exclude

	for end < (d.total - 1) {

		end = start + d.partial - 1

		if end >= d.total {
			end = d.total - 1
		}

		rg := Range{
			Id:    id,
			Start: start,
			End:   end,
		}
		rgs = append(rgs, rg)

		id++
		start = end + 1
	}

	return rgs
}

func rangeGet(url string, start, end uint64) ([]byte, error) {

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

func rangeGets(url string, rs Ranges) ([]byte, error) {

	var bys []byte
	for _, v := range rs {

		by, err := rangeGet(url, v.Start, v.End)
		if err != nil {
			return nil, err
		}
		bys = append(bys, by...)
	}

	return bys, nil
}

type header struct {
	AcceptRanges  bool
	ContentLength uint64
}

func getHeader(url string) (header, error) {

	var h header

	res, err := http.Head(url)
	if err != nil {
		logger.Error(err)
		return h, err
	}

	if n := res.Header.Get("Content-Length"); n != "" {
		i, err := strconv.Atoi(n)
		if err != nil {
			logger.Error(err)
			return h, err
		}
		h.ContentLength = uint64(i)
	}

	if res.Header.Get("Accept-Ranges") == "bytes" {
		h.AcceptRanges = true
	}

	return h, nil

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

type Client struct {
	Size uint64
}

func NewClient(size uint64) *Client {
	return &Client{Size: size}
}

func (rec Client) Get(url string) ([]byte, error) {

	t1 := time.Now()
	h, err := getHeader(url)
	if err != nil {
		return nil, err
	}
	if !h.AcceptRanges {
		return get(url)
	}

	if h.ContentLength <= rec.Size {
		return get(url)
	}

	rs := getRanges(division{h.ContentLength, rec.Size, 0, 0})

	by, err := rangeGets(url, rs)
	if err != nil {
		return nil, err
	}
	diff := time.Now().Sub(t1)
	logger.Debug(diff)
	return by, nil
}
