package server

import (
	"cloud/lib/logger"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"test/grpc/proto/hello"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type HelloServer struct {
	reply map[int]chan string
	port  int
	n     int
	mu    *sync.Mutex
}

func (rec *HelloServer) SayHello(stream hello.HelloService_SayHelloServer) error {

	i, ch := rec.Add()

	waitc := make(chan error, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {

		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop receiving from Deploy_GetMsgClient:", i)
				return
			default:
				_, err := stream.Recv()
				if err != nil {
					logger.Error(err)
					logger.Warn("connection from Deploy_GetMsgClient:", i, "closed")
					waitc <- err
					return
				}

			}
		}

	}()

	go func() {

		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop sending to Deploy_GetMsgClient:", i)
				return
			case s := <-ch:
				reply := &hello.HelloResponse{Reply: s}
				logger.Debug("send:", i, reply)
				if err := stream.Send(reply); err != nil {
					logger.Error(err)
					waitc <- err
					return
				}
			}
		}

	}()

	return <-waitc

}

func NewHelloServer(port int) *HelloServer {

	ser := &HelloServer{
		reply: make(map[int]chan string),
		port:  port,
		mu:    new(sync.Mutex),
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: time.Second * 10,
		}),
	)
	hello.RegisterHelloServiceServer(grpcServer, ser)

	go grpcServer.Serve(lis)

	return ser
}

func (rec HelloServer) Send(s string) {

	for i := range rec.reply {
		rec.reply[i] <- s
	}
}

func (rec *HelloServer) Add() (int, chan string) {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.n++
	ch := make(chan string)
	rec.reply[rec.n] = ch
	return rec.n, ch
}
