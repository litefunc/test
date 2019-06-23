package docker

import (
	"cloud/lib/logger"
	"fmt"
	"time"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
	"cloud/server/ota/service/docker/input"
)

func (rc Service) CreateDocker(token string, in input.CreateDocker) (Docker, error) {

	var data Docker

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return data, err
	}
	defer tx.Rollback()

	data, err = rc.createDocker(token, tx, in)
	if err != nil {
		return data, err
	}

	return data, tx.Commit()
}

func (rc Service) createDocker(token string, tx model.Tx, in input.CreateDocker) (Docker, error) {

	var data Docker

	pdata, err := rc.portal.Get(token, "bio")
	if err != nil {
		return data, err
	}

	user := pdata.Bio


	return newDocker, nil
}

