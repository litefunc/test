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
	logger.Debug("start")
	logger.Debug(connected("http://clients3.google.com/generate_204"))
	logger.Debug(connected("http://clients3.google.com/generate_404"))
	logger.Debug(connected("http://fake"))
	logger.Debug(connected("http://localhost:8080"))
}
