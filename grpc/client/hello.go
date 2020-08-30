package client

import (
	"cloud/lib/logger"
	"context"
	"fmt"
	"net"
	"net/http"
	"test/grpc/proto/hello"
	"time"

	"google.golang.org/grpc"
)

type connDialer struct {
	c net.Conn
}

func (cd connDialer) Dial(network, addr string) (net.Conn, error) {
	return cd.c, nil
}

func NewHttpClient(localPort int, remoteIP string, remotePort int) (*http.Client, error) {

	localIP := []byte{} //  any IP，不指定IP

	netAddr := &net.TCPAddr{Port: localPort}
	if len(localIP) != 0 {
		netAddr.IP = localIP
	}

	fmt.Println("netAddr:", netAddr)

	RemoteEP := net.TCPAddr{IP: net.ParseIP(remoteIP), Port: remotePort}
	conn, err := net.DialTCP("tcp", netAddr, &RemoteEP)
	// conn, err := dialer.Dial("tcp", "localhost:8090")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	conn.SetLinger(0)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: connDialer{conn}.Dial,
		},
	}

	return client, nil

}

type HelloClient struct {
	serverAddr string
	greeting   chan string
	HttpClient *http.Client
}

// func (rec HelloClient) Dial() (*grpc.ClientConn, error) {
// 	kacp := keepalive.ClientParameters{
// 		Time:    100000 * time.Second, // send pings every 10 seconds if there is no activity
// 		Timeout: 100000 * time.Second, // wait 10 second for ping ack before considering the connection dead
// 		// PermitWithoutStream: true,             // send pings even without active streams
// 	}
// 	conn, err := grpc.Dial(rec.serverAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}
// 	co1 := option.WithGRPCConn(conn)
// 	// co2 := option.WithHTTPClient(rec.HttpClient)
// 	conn1, err := transport.Dial(context.Background(), co1)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}
// 	return conn1, err
// }

func (rec HelloClient) Dial() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(rec.serverAddr, grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			dst, err := net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			src := &net.TCPAddr{
				// IP:   net.ParseIP("127.0.0.1"),
				Port: 8091,
			}
			return net.DialTCP("tcp", src, dst)
		}))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return conn, err
}

func (rec *HelloClient) SayHello() {

	sleep := make(chan time.Duration, 1)
	for {
		select {
		case n := <-sleep:
			time.Sleep(time.Second * n)

		default:

			logger.Warnf("Dial:%s", rec.serverAddr)
			// conn, err := grpc.Dial(rec.serverAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
			conn, err := rec.Dial()
			if err != nil {
				sleep <- 5
				conn.Close()
				continue
			}
			defer conn.Close()

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
