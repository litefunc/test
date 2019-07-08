package service

type Service struct {
	Package string
	Model   string
	Models  string
}

func New(pkg, md, mds string) Service {
	return Service{
		Package: pkg,
		Model:   md,
		Models:  mds,
	}
}

const Init = `package {{.Package}}

import (
	"cloud/server/ota/gcs"
	"cloud/server/ota/service/v2/common/model"
	"cloud/server/ota/service/v2/common/portal"
)

type Service struct {
	portal               portal.Agent
	DB                   model.DB
}
`
