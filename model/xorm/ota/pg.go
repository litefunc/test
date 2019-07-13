package model

import (
	"cloud/server/ota/service/v2/common/model"

	"github.com/go-xorm/xorm"
)

type Pg struct {
	*xorm.Engine
}

type Tx struct {
	*xorm.Session
}

func (pg *Pg) Begin() model.Tx {
	tx := pg.NewSession()
	tx.Begin()
	return &Tx{tx}
}

func (pg *Tx) Commit() error {
	return pg.Session.Commit()
}

func (pg *Tx) Rollback() error {
	return pg.Session.Rollback()
}

func NewPg(db *xorm.Engine) *Pg {
	return &Pg{Engine: db}
}

func getTx(tx model.Tx) *Tx {
	return tx.(*Tx)
}
