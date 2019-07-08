package main


import (
	"cloud/lib/logger"
	"test/grpc/localserver"
	"flag"
	"time"
	// "errors"
)

var sn string

func main() {

	flag.StringVar(&sn, "sn", "N0001", "sn")
	flag.Parse()

	var wc chan int

	cli, _ := localserver.NewClient(sn, "ab:01:23:cd:45:67", "localhost:50010", "localhost:50030")
	logger.Debug(cli)

	time.Sleep(time.Second*4)


	<- wc
}