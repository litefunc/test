package deploy

import (
	"cloud/lib/logger"
	"cloud/protos/pb"
	"context"
)

func (ser *DeployServer) Login(ctx context.Context, account *pb.Account) (*pb.Token, error) {
	logger.Info(account)

	return &pb.Token{Token: ""}, nil
}
