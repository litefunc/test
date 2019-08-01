package pgsql

import (
	"test/pgsql/query"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db sqlx.DB
	q  query.Query
}

func (db DB) Select(cols ...string) DB {
	db.q = db.q.Select(cols...)
	return db
}

func (db DB) From(table string) DB {
	db.q = db.q.From(table)
	return db
}

func (db DB) Where(condition string, args ...interface{}) DB {
	db.q = db.q.Where(condition, args...)
	return db
}

func (db DB) Order(args ...string) DB {
	db.q = db.q.Order(args...)
	return db
}

func (db DB) Limit(n uint64) DB {
	db.q = db.q.Limit(n)
	return db
}

func (db DB) SQL() string {
	return db.q.SQL()
}

func (db DB) Args() []interface{} {
	return db.q.Args()
}
