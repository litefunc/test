package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"errors"
	"strings"
	"test/pgsql/query"

	_ "github.com/lib/pq"
)

type Exec struct {
	DB Database
	q  query.Query
	md interface{}
}

func NewExec(db Database, q query.Query, md interface{}) Exec {
	return Exec{DB: db, q: q, md: md}
}

var ErrDeleteWithoutCondition = errors.New(`pgsql: DELETE must have condition. To delete all rows, use TRUNCATE instead`)

func (db Exec) Run() (sql.Result, error) {

	// logger.Debug(db.q.SQL(), db.q.Args())

	if db.q.Type() == query.DELETE {
		if !strings.Contains(db.q.SQL(), "WHERE") {
			return nil, ErrDeleteWithoutCondition
		}
	}

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
