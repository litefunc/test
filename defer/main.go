package main

import (
	"cloud/lib/logger"
	"time"
)

func loop() {
	for i := 0; i < 10; i++ {
		logger.Debug(i)
		defer logger.Debug("defer", i)
	}
	logger.Debug("finish")
}

func main() {
	loop()
	logger.Debug(de())
}

func de() int {
	defer logger.Debug(0)
	return de2(1)
}

func de1(i int) int {
	defer logger.Debug(i + 1)
	logger.Debug(i)
	return i
}

func de2(i int) int {
	defer debug(i + 1)
	logger.Debug(i)
	return i
}

func debug(i int) {

	logger.Debug("sleep")
	time.Sleep(time.Second * 2)
	logger.Debug(i)
}
