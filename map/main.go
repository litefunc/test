package main

import (
	"VodoPlay/logger"
	"encoding/json"
)

func main() {

	var m map[int]int

	um([]byte(`{}`), &m)
	m[1] = 1

}

func j(o interface{}) {
	logger.Debug(o)
	by, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by))
}

func um(data []byte, o interface{}) {
	logger.Debug(string(data))
	logger.Debug(o)

	if err := json.Unmarshal(data, o); err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(o)
}
