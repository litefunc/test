package service

const Get = `package {{.Package}}

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/v1/model"
	"cloud/server/ota/service/common/v1/status"
	"cloud/server/ota/service/common/v1/util"
)

func (rc Service) Get{{.Models}}(token string) ({{.Models}}, error) {

	var data {{.Models}}

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	return rc.get{{.Models}}(token, tx)
}

func (rc Service) get{{.Models}}(token string, tx model.Tx) ({{.Models}}, error) {

	var data {{.Models}}

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return data, err
	}

	user := pdata.Bio

	return data, nil
}

`
