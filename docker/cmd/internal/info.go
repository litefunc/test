package internal

import (
	"cloud/lib/logger"
	"context"

	"github.com/docker/docker/client"
)

func Info() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	cs, err := cli.Info(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(cs.ServerVersion)
}
