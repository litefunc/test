package main

import (
	"fmt"
	"test/grpc/client"
)

func main() {

	host := "localhost"
	port := 8090

	cli := client.NewHelloClient(fmt.Sprintf(`%s:%d`, host, port))
	cli.SayHello()

}
