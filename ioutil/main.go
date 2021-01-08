package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"test/ioutil/internal"
	"test/logger"
)

type A struct {
	A int `json:"a"`
}

func main() {

	var by []byte
	logger.LogErr(ioutil.WriteFile("test.json", by, os.ModePerm))

	var a A
	logger.Debug(by == nil)
	logger.LogErr(json.Unmarshal(by, &a))

	internal.Temp()
}
