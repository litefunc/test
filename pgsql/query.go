package pgsql

import (
	"cloud/lib/logger"
	"test/pgsql/query"

	"github.com/jmoiron/sqlx"
)

type Query struct {
	*sqlx.DB
	q  query.Query
	md interface{}
}

func NewQuery(db *sqlx.DB, q query.Query, md interface{}) Query {
	return Query{DB: db, q: q, md: md}
}

func (db Query) Run() error {
	logger.Debug(db.md)
	logger.Debug(db.q.SQL())
	logger.Debug(db.q.Args()...)

	if err := db.DB.Select(db.md, db.q.SQL(), db.q.Args()...); err != nil {
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

func (db Query) Order(args ...string) Query {
	db.q = db.q.Order(args...)
	return db
}

func (db Query) Limit(n uint64) Query {
	db.q = db.q.Limit(n)
	return db
}
