package model

import (
	"cloud/server/ota/service/v2/common/model"

	"github.com/go-pg/pg"
)

type Pg struct {
	*pg.DB
}

type Tx struct {
	*pg.Tx
}

func (pg *Pg) Begin() (model.Tx, error) {
	return pg.DB.Begin()
}

func (pg *Tx) Commit() error {
	return pg.Tx.Commit()
}

func (pg *Tx) Rollback() error {
	return pg.Tx.Rollback()
}

func NewPg(db *pg.DB) *Pg {
	return &Pg{DB: db}
}

func getTx(tx model.Tx) *Tx {
	return tx.(*Tx)
}
