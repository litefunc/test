package main

import (
	"errors"
	"fmt"
	"test/logger"
)

func main() {

	logger.Debug(fmt.Sprintf("%.2f", 12.345))
	logger.Debug(fmt.Sprintf("%.2f", 12.344))

	logger.Debug(fmt.Sprintf("%s", errors.New("error")))
	logger.Debug(fmt.Sprintf("%s", []byte("byte")))
}
