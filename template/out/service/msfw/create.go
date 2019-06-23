package msfw

import (
	"cloud/lib/logger"
	"fmt"
	"time"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/msfw/input"
)

func (rc Service) CreateMsFw(token string, in input.CreateMsFw) (MsFw, error) {

	var data MsFw

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	data, err = rc.createMsFw(token, tx, in)
	if err != nil {
		return data, err
	}

	return data, tx.Commit()
}

func (rc Service) createMsFw(token string, tx model.Tx, in input.CreateMsFw) (MsFw, error) {

	var data MsFw

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return data, err
	}

	user := pdata.Bio


	return newMsFw, nil
}

