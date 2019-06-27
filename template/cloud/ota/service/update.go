package service

const Update = `package {{.Package}}

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/v1/model"
	"cloud/server/ota/service/common/v1/status"
	"cloud/server/ota/service/v1/{{.Package}}/input"
)

func (rc Service) Update{{.Model}}(token string, in input.Update{{.Model}}) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err = rc.update{{.Model}}(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) update{{.Model}}(token string, tx model.Tx, in input.Update{{.Model}}) error {

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

`
