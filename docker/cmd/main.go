package main

import (
	"cloud/lib/logger"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	cs, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		logger.Error(err)
		return
	}
	for i, v := range cs {
		logger.Debug(i, v)
	}

	m := make(map[uint64]string)
	var n uint64
	for _, v := range cs {
		for _, v := range v.RepoTags {
			n++
			m[n] = v
			logger.Debug(n, v)
		}
	}

}
