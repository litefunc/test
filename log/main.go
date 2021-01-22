package main

import (
	"cloud/lib/logger"
	"log"
	"os"
	"test/log/internal"
)

func main() {
	logger.Debug("start")
	log1 := log.New(os.Stderr, "[DEBUG] ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	log1.Println(0)

	l := internal.NewLogger(log1)
	l.Println(1)
	l.Debug(1)
	l.Info(1)
}
