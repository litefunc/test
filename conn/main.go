package main

import (
	"cloud/lib/logger"
	"io/ioutil"
	"net/http"
)

func connected(url string) (ok bool) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return false
	}
	logger.Debug(url, resp.StatusCode, http.StatusText(resp.StatusCode), resp.Body)
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)
	logger.Debug(string(by))
	return true
}

func main() {
	// logger.Debug("start")

	// logger.Debug(connected("https://www.bing.com"))
	// logger.Debug(connected("https://www.sina.com.cn"))
	// logger.Debug(connected("http://www.baidu.com"))
	// logger.Debug(connected("http://clients3.google.com/generate_204"))
	// logger.Debug(connected("http://clients3.google.com/generate_404"))
	// logger.Debug(connected("http://fake"))
	// logger.Debug(connected("http://localhost:8080"))

	logger.Debug(connected1())
}

func connected1() bool {

	urls := []string{"http://clients3.google.com/generate_204", "https://www.bing.com", "http://www.baidu.com", "https://www.sina.com.cn"}
	for _, v := range urls {
		if conn(v) {
			return true
		}
	}
	return false
}

func conn(url string) (ok bool) {
	_, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func conns(urls ...string) {
	for _, v := range urls {
		logger.Debug(conn(v))
	}
}
