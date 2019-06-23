package routerfw

import (
	"cloud/server/ota/gcs"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/portal"
)

type Service struct {
	portal               portal.Agent
	DB                   model.DB
}
