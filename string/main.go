package main

import (
	"cloud/lib/logger"
	"encoding/json"
	"strings"
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

	s = strings.Replace(" abc  def gh", "  ", " ", -1)
	logger.Debug(s)
	s = strings.Replace(" abc   def gh", "  ", " ", -1)
	logger.Debug(s)

	ss := strings.Split(" abc   def gh", " ")
	for i, s := range ss {
		logger.Debug(i, len(s), s)
	}

	logger.Debug("" >= "2030378@f1b8cd7b")

	var bytes []byte
	bytes = nil
	logger.Debug(string(bytes))
}
