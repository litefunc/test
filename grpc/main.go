package main

import (
	"fmt"
	"test/grpc/client"
	"test/grpc/server"
	"time"
)

func main() {
	port := 8090

	ser := server.NewHelloServer(port)
	time.Sleep(time.Second)

	for i := 0; i < 400; i++ {
		cli := client.NewHelloClient(fmt.Sprintf(`localhost:%d`, port))
		go cli.SayHello()
	}

	time.Sleep(time.Second)
	for _, s := range []string{"a", "b", "c", "d", "e"} {
		ser.Send(s)
	}

	var wc chan struct{}
	<-wc
}
