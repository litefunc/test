package main

import (
	"cloud/lib/logger"
	"log"
	"os"
)

func main() {
	logger.Debug("start")
	log1 := log.New(os.Stderr, "| DEBUG | ", log.Ldate|log.Ltime|log.LUTC|log.Lmsgprefix|log.Lshortfile)
	log1.Println(0)

}
