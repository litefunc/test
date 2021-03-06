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
	defer rec.Delete(i)
	// go rec.tick(i, ch)
	waitc := make(chan error, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {

		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop receiving from:", i)
				return
			default:
				in, err := stream.Recv()
				if err != nil {
					logger.Error(err)
					logger.Warn("connection from:", i, "closed")
					waitc <- err
					return
				}
				logger.Warnf("recv:%+v from:%d", in, i)

			}
		}

	}()

	go func() {

		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop sending to:", i)
				return
			case s := <-ch:
				reply := &hello.HelloResponse{Reply: s}
				logger.Infof("send:%+v to %d", reply, i)
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
			Time:    time.Second * 10,
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
	logger.Info("add:", rec.n)
	return rec.n, ch
}

func (rec *HelloServer) Delete(i int) {
	rec.mu.Lock()
	defer rec.mu.Unlock()

	delete(rec.reply, i)
	logger.Info("delete:", rec.n)
	return
}

func (rec HelloServer) tick(i int, ch chan string) {
	var n int
	for {
		time.Sleep(time.Second * 3)
		ch <- fmt.Sprintf("server msg %d %d", i, n)
		n++
	}
}
