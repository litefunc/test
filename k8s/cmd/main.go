package main

import (
	"cloud/lib/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func get(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
	}
	by, _ := ioutil.ReadAll(resp.Body)
	s := string(by)
	n := strings.Split(s, "visits: ")[1]

	i, err := strconv.Atoi(strings.Replace(n, "\n", "", -1))
	if err != nil {
		logger.Error(err)
	}
	return i
}

func loadbalance(url string) bool {
	n := get(url)
	for i := 0; i < 100; i++ {
		n1 := get(url)
		if n1-n <= 0 {
			logger.Debug(n1, n)
			return true
		}
		n = n1
	}
	return false
}

func testLoadbalance(url string, n int) []bool {
	var wc sync.WaitGroup

	var success []bool
	for i := 0; i < n; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()
			success = append(success, loadbalance(url))

		}(i)
	}
	wc.Wait()

	return success
}

func main() {
	// url := "http://192.168.0.12:31327"
	url := "http://10.96.77.158:8088"
	list := testLoadbalance(url, 3)
	for _, v := range list {
		if v {
			logger.Debug(v)
		} else {
			logger.Error(v)
		}
	}

}
