package main

import (
	"MediaImage/logger"
	"time"
)

func main() {
	logger.Debug("start")
	var n int

	t := time.NewTicker(time.Second * 2)
	for c := range t.C {
		logger.Debug(n, c)
		n++
	}
}
