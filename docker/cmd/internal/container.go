package internal

import (
	"cloud/lib/logger"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ContainerList() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	cs, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		logger.Error(err)
		return
	}
	for i, v := range cs {
		logger.Debug(i, v)
		// c, err := cli.ContainerInspect(ctx, v.ID)
		// if err != nil {
		// 	logger.Error(err)
		// }
		// by, _ := json.Marshal(c)
		// logger.Info(i, string(by))
	}

}
