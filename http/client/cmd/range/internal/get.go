package internal

import (
	"VodoPlay/logger"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"net/http"
	"strconv"
	"time"
)

type Range struct {
	Start uint64
	End   uint64
}

type Ranges []Range

type Division struct {
	Total   uint64
	Partial uint64
	Start   uint64
}

func NewDivision(total, partial, start uint64) Division {
	return Division{Total: total, Partial: partial, Start: start}
}

func (rec Division) Ranges() (Ranges, error) {

	rgs := Ranges{}

	// return if parameter invalid
	if rec.Total == 0 || rec.Partial == 0 || rec.Total < rec.Start {
		return rgs, fmt.Errorf(`Parameter invalid:%+v`, rec)
	}

	start := rec.Start
	end := rec.Start

	for end < (rec.Total - 1) {

		end = start + rec.Partial - 1

		if end >= rec.Total {
			end = rec.Total - 1
		}

		rg := Range{
			Start: start,
			End:   end,
		}
		logger.Debug(rg)
		rgs = append(rgs, rg)

		start = end + 1
	}

	return rgs, nil
}

func (rec Range) GetResponse(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rec.Start, rec.End))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return resp, nil
}

func (rec Range) GetBody(url string) ([]byte, error) {

	resp, err := rec.GetResponse(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return by, nil
}

func (rec Ranges) GetBody(url string) ([]byte, error) {

	var bys []byte
	for _, v := range rec {

		by, err := v.GetBody(url)
		if err != nil {
			return bys, err
		}
		bys = append(bys, by...)
	}

	return bys, nil
}

func (rec Range) DownloadFile(url, filepath string) error {

	resp, err := rec.GetResponse(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (rec Ranges) DownloadFile(url, filepath string) error {

	for _, v := range rec {
		if err := v.DownloadFile(url, filepath); err != nil {
			return err
		}
	}

	return nil
}

func (rec Division) GetBody(url string) ([]byte, error) {

	rs, err := rec.Ranges()
	if err != nil {
		return nil, err
	}

	return rs.GetBody(url)
}

func (rec Division) DownloadFile(url, filepath string) error {

	rs, err := rec.Ranges()
	if err != nil {
		return err
	}

	return rs.DownloadFile(url, filepath)
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

func Get(url string) ([]byte, error) {
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
	logger.Debug(len(by))

	diff := time.Now().Sub(t1)
	logger.Debug(diff)

	h := md5.New()
	if _, err := h.Write(by); err != nil {
		logger.Error(err)
		return nil, err
	}
	ck := fmt.Sprintf(`%x`, h.Sum(nil))
	logger.Debug(ck)

	return by, nil
}

func DownloadFile(url, filepath string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer f.Close()

	// Write the body to file
	_, err = io.Copy(f, resp.Body)
	return err
}

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (rec Client) DownloadFile(url, filepath string) error {
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		logger.Error(err)
		return err
	}

	t1 := time.Now()

	if err := DownloadFile(url, filepath); err != nil {
		return err
	}

	diff := time.Now().Sub(t1)
	logger.Debug(diff)

	s, err := Md5Sum(filepath)
	if err != nil {
		return err
	}
	logger.Debug(s)

	return nil
}

func (rec Client) RangeDownloadFile(url, filepath string, d Division) error {
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		logger.Error(err)
		return err
	}

	t1 := time.Now()

	if err := d.DownloadFile(url, filepath); err != nil {
		return err
	}

	diff := time.Now().Sub(t1)
	logger.Debug(diff)

	s, err := Md5Sum(filepath)
	if err != nil {
		return err
	}
	logger.Debug(s)

	return nil
}

func (rec Client) Get(url string, partial uint64) ([]byte, error) {

	t1 := time.Now()
	h, err := getHeader(url)
	if err != nil {
		return nil, err
	}
	if !h.AcceptRanges {
		return Get(url)
	}

	if h.ContentLength <= partial {
		return Get(url)
	}

	rs, err := Division{h.ContentLength, partial, 0}.Ranges()

	by, err := rs.GetBody(url)
	if err != nil {
		return nil, err
	}
	diff := time.Now().Sub(t1)
	logger.Debug(diff)
	return by, nil
}

func Md5Sum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		logger.Error(err)
		return "", err
	}

	return fmt.Sprintf(`%x`, h.Sum(nil)), nil
}
