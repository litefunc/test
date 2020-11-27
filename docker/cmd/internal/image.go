package internal

import (
	"cloud/lib/logger"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ImageList() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	list, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		logger.Error(err)
		return
	}
	for i, v := range list {
		logger.Debug(i, v)
	}

}
