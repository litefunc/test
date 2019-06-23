package msfw

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/common/util"
)

func (rc Service) GetMsFws(token string) (MsFws, error) {

	var data MsFws

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	return rc.getMsFws(token, tx)
}

func (rc Service) getMsFws(token string, tx model.Tx) (MsFws, error) {

	var data MsFws

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


