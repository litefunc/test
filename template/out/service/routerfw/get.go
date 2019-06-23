package routerfw

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/common/util"
)

func (rc Service) GetRouterFws(token string) (RouterFws, error) {

	var data RouterFws

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	return rc.getRouterFws(token, tx)
}

func (rc Service) getRouterFws(token string, tx model.Tx) (RouterFws, error) {

	var data RouterFws

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


