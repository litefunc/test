package routerfw

import (
	"cloud/lib/logger"
	"fmt"
	"time"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/routerfw/input"
)

func (rc Service) CreateRouterFw(token string, in input.CreateRouterFw) (RouterFw, error) {

	var data RouterFw

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	data, err = rc.createRouterFw(token, tx, in)
	if err != nil {
		return data, err
	}

	return data, tx.Commit()
}

func (rc Service) createRouterFw(token string, tx model.Tx, in input.CreateRouterFw) (RouterFw, error) {

	var data RouterFw

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return data, err
	}

	user := pdata.Bio


	return newRouterFw, nil
}

