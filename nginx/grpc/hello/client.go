package hello

import (
	"cloud/lib/logger"
	"context"
	"fmt"
	"test/nginx/pb"

	"google.golang.org/grpc"
)

type client struct {
	Target string
}

func NewClient(host string, port int) *client {
	cli := client{Target: fmt.Sprintf("%s:%d", host, port)}
	return &cli
}

func (cli client) Hello(report *pb.Empty) (*pb.Reply, error) {

	conn, err := grpc.Dial(cli.Target, grpc.WithInsecure())
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewHelloClient(conn)
	reply, err := client.Hello(context.Background(), report)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return reply, nil
}
