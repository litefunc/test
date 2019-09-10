package main

import (
	"VodoPlay/logger"
	"net/http"
)

func connected(url string) (ok bool) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return false
	}
	logger.Debug(resp.StatusCode)
	return true
}

func main() {
	logger.Debug(connected("http://clients3.google.com/generate_204"))
	logger.Debug(connected("http://clients3.google.com/generate_404"))
	logger.Debug(connected("http://fake"))
}
