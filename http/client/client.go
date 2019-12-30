package main

import (
	"cloud/lib/logger"
	"io/ioutil"
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

func get2(i int) {
	resp, err := http.Get("http://localhost:8080/sleep100")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(i, resp.StatusCode)
	defer resp.Body.Close()
}

func get3(i int) {
	resp, err := http.Get("http://localhost:8080/index")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(i, resp.StatusCode)
	logger.Debug(resp.Body.Close())
	logger.Debug(resp.Body.Close())

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by))
}

func connected(url string) (ok bool) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return false
	}
	logger.Debug(resp.StatusCode)
	return true
}

func connected1(url string) (ok bool) {
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
	var wc sync.WaitGroup
	wc.Add(1)
	for i := 0; i < 10; i++ {
		// get3(i)
		connected("http://clients3.google.com/generate_204")
	}
	wc.Wait()
}
