package hello

import (
	"cloud/lib/logger"
	"context"

	"fmt"
	"net"

	"test/nginx/pb"

	"google.golang.org/grpc"
)

type Server struct {
	id uint64
}

func NewServer(id uint64, port int) *Server {

	addr := fmt.Sprintf(`:%d`, port)
	lis, err := net.Listen("tcp", addr)

	logger.Infof("Server listen at %v", addr)
	if err != nil {
		logger.Fatal("failed to listen:", err)
	}
	grpcServer := grpc.NewServer()

	ser := Server{id}

	pb.RegisterHelloServer(grpcServer, &ser)

	go grpcServer.Serve(lis)

	return &ser
}

func (s *Server) Hello(ctx context.Context, in *pb.Empty) (*pb.Reply, error) {
	return &pb.Reply{Id: s.id}, nil
}
