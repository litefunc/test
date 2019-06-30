package service

const Get = `package {{.Package}}

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/v1/common/model"
	"cloud/server/ota/service/v1/common/status"
	"cloud/server/ota/service/v1/common/util"
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

func (rc Service) Get{{.Models}}ByGroup(token string, gp uint64) ({{.Models}}, error) {
	var data {{.Models}}

	pdata, err := rc.portal.Get(token, "groups")
	if err != nil {
		return data, err
	}

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	if !util.IdValid(pdata.Groups.IDs(), gp) {
		return data, status.ErrCheckGroup(gp)
	}

	return rc.{{.Model}}Store.{{.Models}}ByGroup(tx, gp)
}

func (rc Service) Get{{.Models}}ByCompany(token string, com uint64) ({{.Models}}, error) {
	var data {{.Models}}

	pdata, err := rc.portal.Get(token, "coms")
	if err != nil {
		return data, err
	}

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	if !util.IdValid(pdata.Companies.IDs(), com) {
		return data, status.ErrCheckCompany(com)
	}

	return rc.{{.Model}}Store.{{.Models}}ByCom(tx, com)
}


`
