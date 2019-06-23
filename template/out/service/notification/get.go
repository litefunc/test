package notification

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/common/util"
)

func (rc Service) GetNotifications(token string) (Notifications, error) {

	var data Notifications

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	return rc.getNotifications(token, tx)
}

func (rc Service) getNotifications(token string, tx model.Tx) (Notifications, error) {

	var data Notifications

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

