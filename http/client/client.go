package main

import (
	"cloud/lib/logger"
	"net/http"
	"sync"
)

func sleep() {
	var wc sync.WaitGroup
	for i := 0; i < 10; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()

			_, err := http.Get("http://localhost:8080/sleep")
			if err != nil {
				logger.Error(err)
			}

		}(i)
	}
	wc.Wait()

	for i := 0; i < 10; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()

			_, err := http.Get("http://localhost:8080/sleep1")
			if err != nil {
				logger.Error(err)
			}

		}(i)
	}
	wc.Wait()
}

func get(i int) {
	resp, err := http.Get("http://localhost:8080/index")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(i, resp.StatusCode)
}

func get1(i int) {
	resp, err := http.Get("http://localhost:8080/index")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(i, resp.StatusCode)
	defer resp.Body.Close()
}

func main() {
	var wc sync.WaitGroup
	wc.Add(1)
	for i := 0; i < 10; i++ {
		get(i)
	}
	wc.Wait()
}
