package deploy

import (
	"cloud/lib/logger"
	"cloud/protos/pb"
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

type DevMsg struct {
	Sn   string
	To   string
	Data []byte
}

func (s *DeployServer) GetMsg(stream pb.Deploy_GetMsgServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		logger.Error("unable to get client's metadata")
		return errors.New("metadata is required")
	}
	sns, ok := md["sn"]
	if !ok || len(sns) == 0 || sns[0] == "" {
		logger.Error("unable to get client's sn")
		return  errors.New("please provide sn in metadata")
	}
	sn := sns[0]

	cl := MsgClient{Id: sn, Ch: make(chan pb.Msg)}
	s.MsgServer.ClientIn <- cl
	defer s.MsgServer.Delete(cl)


	waitc := make(chan error, 1)
	s.Waitc[sn] = waitc

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop receiving from Deploy_GetMsgClient:", sn)
				return
			default:
				devMsg, err := stream.Recv()
				if err == nil {
				
				}
				if err != nil {
					logger.Error(sn, err)
					waitc <- err
					return
				}
				logger.Debug(devMsg)

			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Warn("stop sending to Deploy_GetMsgClient:", sn)
				return
			case msg := <-cl.Ch:
				logger.Infof(`Deploy_GetMsg Service Send => sn:%s, Msg: {from:%s, update:%+v}`, sn, msg.From, msg.Update)
				if err := stream.Send(&msg); err != nil {
					logger.Error(err)
					waitc <- err
					return
				}
			}
		}
	}()

	err := <-waitc
	logger.Warn(sn , "return")
	return err
}


func (ser *DeployServer) listen() {
	ser.Noti = make(chan pb.Noti)
	ser.DevMsg = make(chan DevMsg)

	go func() {
		for {
			noti := <-ser.Noti
			if noti.Msg != nil {
				msg := pb.Msg{From: noti.Msg.From, Update: noti.Msg.Update}
				for _, sn := range noti.SnList {
					ser.MsgServer.MsgCh <- MsgTo{Id: sn, Msg: msg}
				}
			}
		}
	}()

	go func() {
		for {
			devMsg := <-ser.DevMsg

			deployMsg := pb.DeplMsg{Sn: devMsg.Sn, Data: devMsg.Data}

			noti, err := ser.DeplMsg[devMsg.To].GetNoti(&deployMsg)
			if err != nil || noti == nil {
				continue
			}
			ser.Noti <- *noti
		}
	}()
}