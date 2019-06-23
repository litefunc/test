package docker

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/common/util"
)

func (rc Service) GetDockers(token string) (Dockers, error) {

	var data Dockers

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	return rc.getDockers(token, tx)
}

func (rc Service) getDockers(token string, tx model.Tx) (Dockers, error) {

	var data Dockers

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return data, err
	}

	user := pdata.Bio

	return data, nil
}

ser := pdata.Bio

	return data, nil
}

