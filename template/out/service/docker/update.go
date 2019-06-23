package docker

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/docker/input"
)

func (rc Service) UpdateDocker(token string, in input.UpdateDocker) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err = rc.updateDocker(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) updateDocker(token string, tx model.Tx, in input.UpdateDocker) error {

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

