package msfw

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/msfw/input"
)

func (rc Service) UpdateMsFw(token string, in input.UpdateMsFw) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err = rc.updateMsFw(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) updateMsFw(token string, tx model.Tx, in input.UpdateMsFw) error {

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

