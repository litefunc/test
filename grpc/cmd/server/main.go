package main

import (
	"test/grpc/server"
)

func main() {
	port := 8090

	ser := server.NewHelloServer(port)

	for _, s := range []string{"a", "b", "c"} {
		ser.Send(s)
	}

	var wc chan struct{}
	<-wc
}
