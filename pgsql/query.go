package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"reflect"
	"test/pgsql/query"

	"github.com/jmoiron/sqlx"
)

type Query struct {
	DB Database
	q  query.Query
	md interface{}
}

type Database interface {
	Select(dest interface{}, query string, args ...interface{}) error
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func NewQuery(db Database, q query.Query, md interface{}) Query {
	return Query{DB: db, q: q, md: md}
}

func (db Query) Run() error {
	// logger.Debug(db.q.SQL(), db.q.Args())

	t := reflect.TypeOf(db.md)

	var slice bool
	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() == reflect.Slice {
			slice = true
		}
	}

	if t.Kind() == reflect.Slice {
		slice = true
	}

	if slice {
		if err := db.DB.Select(db.md, db.q.SQL(), db.q.Args()...); err != nil {
			logger.Error(err)
			return err
		}
		return nil
	}

	if err := db.DB.QueryRowx(db.q.SQL(), db.q.Args()...).StructScan(db.md); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (db Query) SQL() string {
	return db.q.SQL()
}

func (db Query) Args() []interface{} {
	return db.q.Args()
}

func (db Query) Where(condition string, args ...interface{}) Query {
	db.q = db.q.Where(condition, args...)
	return db
}

func (db Query) GroupBy(cols ...string) Query {
	db.q = db.q.GroupBy(cols...)
	return db
}

func (db Query) Having(condition string, args ...interface{}) Query {
	db.q = db.q.Having(condition, args...)
	return db
}

func (db Query) OrderBy(args ...string) Query {
	db.q = db.q.OrderBy(args...)
	return db
}

func (db Query) Limit(n uint64) Query {
	db.q = db.q.Limit(n)
	return db
}

func (db Query) Offset(n uint64) Query {
	db.q = db.q.Offset(n)
	return db
}
