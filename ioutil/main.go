package main

import (
	"cloud/lib/logger"
	"encoding/json"
	"io/ioutil"
)

type A struct {
	A int `json:"a"`
}

func main() {

	var by []byte
	if err := ioutil.WriteFile("test.json", by, 0755); err != nil {
		logger.Error(err)
	}

	var a A
	logger.Debug(by == nil)
	if err := json.Unmarshal(by, &a); err != nil {
		logger.Error(err)
	}

}
