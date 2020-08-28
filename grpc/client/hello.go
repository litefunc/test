package client

import (
	"cloud/lib/logger"
	"context"
	"fmt"
	"test/grpc/proto/hello"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type HelloClient struct {
	serverAddr string
	greeting   chan string
}

func (rec *HelloClient) SayHello() {

	kacp := keepalive.ClientParameters{
		Time:    100000 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout: 100000 * time.Second, // wait 10 second for ping ack before considering the connection dead
		// PermitWithoutStream: true,             // send pings even without active streams
	}

	sleep := make(chan time.Duration, 1)
	for {
		select {
		case n := <-sleep:
			time.Sleep(time.Second * n)

		default:

			logger.Warnf("Dial:%s", rec.serverAddr)
			conn, err := grpc.Dial(rec.serverAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
			if err != nil {
				sleep <- 5
				conn.Close()
				continue
			}

			logger.Warn("SayHello")
			if err := rec.sayHello(conn); err != nil {
				conn.Close()
				sleep <- 5
				continue
			}
			conn.Close()
			logger.Warn("connection close")
		}
	}

}

func (rec *HelloClient) sayHello(conn *grpc.ClientConn) error {

	client := hello.NewHelloServiceClient(conn)

	waitc := make(chan error, 2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.SayHello(ctx)
	if err != nil {
		logger.Error(err)
		return err
	}

	go func() {
		var n int

		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop receiving")
				return
			default:
				in, err := stream.Recv()
				if err != nil {
					logger.Error(err)
					waitc <- err
					return
				}
				logger.Warn("recv:", in, n)
				n++
			}
		}

	}()

	go func() {
		var n int

		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop sending")
				return
			case s := <-rec.greeting:

				req := &hello.HelloRequest{Greeting: s}
				logger.Info("send:", req, n)
				if err := stream.Send(req); err != nil {
					logger.Error(err)
					waitc <- err
					return
				}
				n++
			}
		}

	}()

	return <-waitc
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
	// go cli.Send()

	return cli
}
