package main

import (
	"MediaImage/logger"
	"log"
	"os"
)

func main() {

	log.Println(1)
	logger.Debug(2)
	logger.Debug(3)
	log.Println(4)

	panic(10)
	os.Exit(1)
}
