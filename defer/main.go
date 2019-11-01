package main

import (
	"cloud/lib/logger"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		logger.Debug(i)
		defer logger.Debug("defer", i)
	}
	logger.Debug("finish")
	time.Sleep(time.Second * 5)
}
