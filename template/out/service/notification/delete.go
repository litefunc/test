package notification

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
)

func (rc Service) DeleteNotification(token string, in NotificationUnique) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err := rc.deleteNotification(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) deleteNotification(token string, tx model.Tx, in NotificationUnique) error {

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

turn data, nil
}

