package main

import (
	"cloud/lib/logger"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.google.com/")
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(resp.StatusCode)
}
