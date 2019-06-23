package service

const Delete = `package {{.Package}}

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/common/model"
	"cloud/server/ota/service/common/status"
)

func (rc Service) Delete{{.Model}}(token string, in {{.Model}}Unique) error {

	tx, err := rc.DB.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	if err := rc.delete{{.Model}}(token, tx, in); err != nil {
		return err
	}

	return tx.Commit()
}

func (rc Service) delete{{.Model}}(token string, tx model.Tx, in {{.Model}}Unique) error {

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
