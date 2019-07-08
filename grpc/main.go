package main


import (
	"cloud/lib/logger"
	"test/grpc/deploy"
	"time"
)


func main() {

	var wc chan int

	ser := deploy.NewDeployServer(50010)
	logger.Debug(ser)

	time.Sleep(time.Second*5)
	ser.Waitc["N0002"] <- nil

	<- wc
}