package main

import (
	"bytes"
	"cloud/lib/logger"
)

func p() {
	panic(1)
}

func r() {
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("Recovered in f", r)
		}
	}()

	p()
}

func main() {

	r()
	var loadBuf bytes.Buffer

	logger.Debug(loadBuf.String())
}
