package client

import (
	"cloud/lib/logger"
	"context"
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
		for {
			in, err := stream.Recv()
			if err != nil {
				logger.Error(err)
				continue
			}
			logger.Debug("receive:", in)
			rec.greeting <- in.Reply
			// logger.Debug("receive 1:", in)
		}
	}()

	go func() {
		// for s := range rec.greeting {
		// 	req := &hello.HelloRequest{Greeting: s}
		// 	if err := stream.Send(req); err != nil {
		// 		logger.Error(err)
		// 	}
		// }
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

func NewHelloClient(serverAddr string) *HelloClient {
	cli := &HelloClient{
		serverAddr: serverAddr,
		greeting:   make(chan string),
	}

	go cli.Listen()

	return cli
}
