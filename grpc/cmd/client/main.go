package main

import (
	"fmt"
	"test/grpc/client"
	"time"
)

func main() {

	host := "localhost"
	port := 8090

	cli := client.NewHelloClient(fmt.Sprintf(`%s:%d`, host, port))
	// httpcli, err := client.NewHttpClient(8091, host, 8090)
	// if err != nil {
	// 	return
	// }
	// cli.HttpClient = httpcli
	for i := 0; i < 10; i++ {
		cli.Dial()
	}
	// cli.SayHello()
	time.Sleep(time.Second * 10)

}
