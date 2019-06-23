package notification

import (
	"cloud/lib/logger"
	"fmt"
	"time"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/notification/input"
)

func (rc Service) CreateNotification(token string, in input.CreateNotification) (Notification, error) {

	var data Notification

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	data, err = rc.createNotification(token, tx, in)
	if err != nil {
		return data, err
	}

	return data, tx.Commit()
}

func (rc Service) createNotification(token string, tx model.Tx, in input.CreateNotification) (Notification, error) {

	var data Notification

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return data, err
	}

	user := pdata.Bio


	return newNotification, nil
}

