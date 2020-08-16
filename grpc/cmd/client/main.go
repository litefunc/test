package main

import (
	"fmt"
	"test/grpc/client"
)

func main() {
	port := 8090

	cli := client.NewHelloClient(fmt.Sprintf(`localhost:%d`, port))
	go cli.SayHello()

	var wc chan struct{}
	<-wc
}
