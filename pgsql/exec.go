package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"test/pgsql/query"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	logger.Debug(db.q.SQL(), db.q.Args())

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

func (db Exec) Set(condition string, args ...interface{}) Exec {
	db.q = db.q.Set(condition, args...)
	return db
}

func (db Exec) Returning(cols ...string) QueryRow {
	db.q = db.q.Returning(cols...)
	return NewQueryRow(db.DB, db.q)
}
