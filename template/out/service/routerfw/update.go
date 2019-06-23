package routerfw

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/routerfw/input"
)

func (rc Service) UpdateRouterFw(token string, in input.UpdateRouterFw) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err = rc.updateRouterFw(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) updateRouterFw(token string, tx model.Tx, in input.UpdateRouterFw) error {

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

