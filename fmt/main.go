package main

import (
	"VodoPlay/logger"
	"fmt"
)

func main() {

	logger.Debug(fmt.Sprintf("%.2f", 12.345))
	logger.Debug(fmt.Sprintf("%.2f", 12.344))
}
