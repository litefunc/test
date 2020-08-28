package main

import (
	"VodoPlay/logger"
	"net/http"
	"time"
)

func main() {

	for {
		conn()
		time.Sleep(time.Second * 5)
	}
}

func conn() {
	res, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		logger.Error(err)
		return
	}
	defer res.Body.Close()
}
