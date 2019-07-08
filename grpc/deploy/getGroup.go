package deploy

import (
	"cloud/protos/pb"
	"context"
)

func (ser *DeployServer) GetGroupId(ctx context.Context, null *pb.Null) (*pb.GroupId, error) {

	return &pb.GroupId{Gid: 0}, nil
}
