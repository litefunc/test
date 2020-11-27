package internal

import (
	"cloud/lib/logger"
	"context"

	"github.com/docker/docker/client"
)

func Ping() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	cs, err := cli.Ping(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(cs)
}
