package main

import (
	"VodoPlay/logger"
	"net"
	"time"
)

func main() {
	port := "8600"

	_, err := net.Listen("tcp", ":"+port)

	if err != nil {
		logger.Fatal(err)
		return
	}

	// if err := ln.Close(); err != nil {
	// 	logger.Fatal(err)
	// 	return
	// }
	time.Sleep(time.Second * 5)
}
