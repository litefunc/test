package deploy

import (
	"cloud/lib/logger"
	"cloud/protos/pb"
	"flag"
	"fmt"
	"net"
	"cloud/server/deploy/grpc/service/deploy"
	"google.golang.org/grpc"
)

type DeployServer struct {
	deploy.DeployServer
	MsgServer *MsgServer
	Noti         chan pb.Noti
	DevMsg       chan DevMsg
	Waitc map[string]chan error
}

func NewDeployServer(port int) *DeployServer {
	flag.Parse()
	addr := fmt.Sprintf(`:%d`, port)
	lis, err := net.Listen("tcp", addr)

	logger.Infof("DeployServer listen at %v", addr)
	if err != nil {
		logger.Fatal("failed to listen:", err)
	}
	grpcServer := grpc.NewServer()

	ser := DeployServer{
		MsgServer:    NewMsgServer("GetMsg"),
		Waitc: make(map[string]chan error),
	}

	pb.RegisterDeployServer(grpcServer, &ser)

	go grpcServer.Serve(lis)

	ser.listen()

	return &ser
}



type MsgServer struct {
	ClientIn   chan MsgClient
	ClientOut  chan MsgClient
	ClientPool map[string]chan pb.Msg
	MsgCh      chan MsgTo
}

type MsgClient struct {
	Id string
	Ch chan pb.Msg
}

type MsgTo struct {
	Id  string
	Msg pb.Msg
}

func NewMsgServer(name string) *MsgServer {
	var ser MsgServer
	ser.MsgCh = make(chan MsgTo)
	ser.ClientPool = make(map[string]chan pb.Msg)
	ser.ClientIn = make(chan MsgClient)
	ser.ClientOut = make(chan MsgClient)

	go func() {
		for {
			select {
			case cl := <-ser.ClientIn:
				ser.ClientPool[cl.Id] = cl.Ch
				logger.Infof(`add %s. %d %sClients registerd`, cl.Id, len(ser.ClientPool), name)

			case cl := <-ser.ClientOut:
				delete(ser.ClientPool, cl.Id)
				logger.Infof(`remove %s. %d %sClients registerd`, cl.Id, len(ser.ClientPool), name)
			case mt := <-ser.MsgCh:
				if ch, ok := ser.ClientPool[mt.Id]; ok {
					ch <- mt.Msg
				}
			}
		}
	}()

	return &ser
}

func (s *MsgServer) Delete(cl MsgClient) {
	s.ClientOut <- cl
}
