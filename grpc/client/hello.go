package client

import (
	"cloud/lib/logger"
	"context"
	"fmt"
	"test/grpc/proto/hello"
	"time"

	"google.golang.org/grpc"
)

type HelloClient struct {
	serverAddr string
	greeting   chan string
}

func (rec *HelloClient) SayHello() {

	conn, err := grpc.Dial(rec.serverAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error(err)
		return
	}
	defer conn.Close()

	client := hello.NewHelloServiceClient(conn)

	stream, err := client.SayHello(context.Background())

	go func() {
		var n int
		for {
			in, err := stream.Recv()
			if err != nil {
				logger.Error(err)
				return
			}
			logger.Warn("recv:", in, n)
			n++

		}
	}()

	go func() {
		var n int
		for s := range rec.greeting {
			req := &hello.HelloRequest{Greeting: s}
			logger.Info("send:", req, n)
			if err := stream.Send(req); err != nil {
				logger.Error(err)
			}
			n++
		}
	}()
	var wc chan struct{}
	<-wc
}

func (rec *HelloClient) Listen() {

	for s := range rec.greeting {
		logger.Debug(s)
		time.Sleep(time.Second * 3)
	}
}

func (rec *HelloClient) Send() {
	var n int
	for {
		rec.greeting <- fmt.Sprintf(`client msg %d`, n)
		time.Sleep(time.Second * 3)
		n++
	}
}

func NewHelloClient(serverAddr string) *HelloClient {
	cli := &HelloClient{
		serverAddr: serverAddr,
		greeting:   make(chan string),
	}

	// go cli.Listen()
	go cli.Send()

	return cli
}
