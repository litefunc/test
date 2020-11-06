package main

import (
	"cloud/lib/logger"
	"io/ioutil"
	"net/http"
)

func get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(url, resp.StatusCode, http.StatusText(resp.StatusCode), resp.Body)
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)

	for i, v := range by {
		logger.Debug(i, v, string(by[i]))
	}
	logger.Debug(len(by))

}

func main() {

	// get("http://localhost:8081/static/test a b .json")
	// get("http://localhost:8081/static/test%20a%20b%20.json")

	// get("http://localhost:8081/static/test%20%E4%B8%AD%E6%96%87%20.txt")
	// get("http://localhost:8081/static/test 中文 .txt")
	// get("http://localhost:8081/static/test%20中文%20.txt")

	get("")
}
