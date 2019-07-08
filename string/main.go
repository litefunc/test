package main

import (
	"cloud/lib/logger"
	"encoding/json"
)

func main() {
	s := "true"
	by, err := json.Marshal(s)
	if err != nil {
		logger.Error(err)
	}
	var b bool
	if err := json.Unmarshal(by, &b); err != nil {
		logger.Error(err)
	}
	logger.Debug(b)

	logger.Debug(string(by))
	if err := json.Unmarshal([]byte(string(by)), &b); err != nil {
		logger.Error(err)
	}

	logger.Debug(b)
}
