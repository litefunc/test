package notification

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/notification/input"
)

func (rc Service) UpdateNotification(token string, in input.UpdateNotification) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err = rc.updateNotification(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) updateNotification(token string, tx model.Tx, in input.UpdateNotification) error {

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return err
	}

	user := pdata.Bio

	if !user.Admin {
		return status.ErrNotAdmin()
	}


	return nil

}

