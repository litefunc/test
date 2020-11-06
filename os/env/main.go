package main

import (
	"VodoPlay/logger"
	"os"
)

func main() {

	logger.Debug(os.Getenv("TESTENV"))
	os.Setenv("TESTENV", "a")
	logger.Debug(os.Getenv("TESTENV"))
}
