package msfw

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
)

func (rc Service) DeleteMsFw(token string, in MsFwUnique) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err := rc.deleteMsFw(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) deleteMsFw(token string, tx model.Tx, in MsFwUnique) error {

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



