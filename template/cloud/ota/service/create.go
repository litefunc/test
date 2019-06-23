package service

const Create = `package {{.Package}}

import (
	"cloud/lib/logger"
	"fmt"
	"time"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/{{.Package}}/input"
)

func (rc Service) Create{{.Model}}(token string, in input.Create{{.Model}}) ({{.Model}}, error) {

	var data {{.Model}}

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	data, err = rc.create{{.Model}}(token, tx, in)
	if err != nil {
		return data, err
	}

	return data, tx.Commit()
}

func (rc Service) create{{.Model}}(token string, tx model.Tx, in input.Create{{.Model}}) ({{.Model}}, error) {

	var data {{.Model}}

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return data, err
	}

	user := pdata.Bio


	return new{{.Model}}, nil
}

`
