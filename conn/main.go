package main

import (
	"VodoPlay/logger"
	"net/http"
)

func connected() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func main() {
	logger.Debug(connected())
}
