package localserver

import (
	"LocalServer/messenger/deploy"
	"LocalServer/messenger/pb"
)

func NewClient(sn, mac, deployAddr, statusAddr string) (*deploy.Client, deploy.MsgCh){

	cli, _ := deploy.NewClient(sn, mac, deployAddr, statusAddr)

	msgc := deploy.MsgCh{Send: nil, Recv: make(chan pb.DevMsg, 1)}
	go deploy.CheckDeploy(cli, deployAddr, msgc)

	return cli, msgc
}


