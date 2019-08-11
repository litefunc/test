package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"test/pgsql/query"

	"github.com/jmoiron/sqlx"
)

type Exec struct {
	*sqlx.DB
	q  query.Query
	md interface{}
}

func NewExec(db *sqlx.DB, q query.Query, md interface{}) Exec {
	return Exec{DB: db, q: q, md: md}
}

func (db Exec) Run() (sql.Result, error) {
	logger.Debug(db.md)
	logger.Debug(db.q.SQL())
	logger.Debug(db.q.Args()...)

	result, err := db.DB.Exec(db.q.SQL(), db.q.Args()...)
	if err != nil {
		logger.Error(err)
		return result, err
	}

	return result, nil
}

func (db Exec) SQL() string {
	return db.q.SQL()
}

func (db Exec) Args() []interface{} {
	return db.q.Args()
}

func (db Exec) Where(condition string, args ...interface{}) Exec {
	db.q = db.q.Where(condition, args...)
	return db
}

func (db Exec) Returning(cols ...string) Exec {
	db.q = db.q.Returning(cols...)
	return db
}
