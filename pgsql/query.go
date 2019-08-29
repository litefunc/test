package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"reflect"
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
	logger.Debug(db.q.SQL(), db.q.Args())

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

func (db Query) Order(args ...string) Query {
	db.q = db.q.Order(args...)
	return db
}

func (db Query) Limit(n uint64) Query {
	db.q = db.q.Limit(n)
	return db
}

type QueryRow struct {
	*sqlx.DB
	q query.Query
}

func NewQueryRow(db *sqlx.DB, q query.Query) QueryRow {
	return QueryRow{DB: db, q: q}
}

func (db QueryRow) Run() *sql.Row {
	logger.Debug(db.q.SQL(), db.q.Args())
	return db.DB.QueryRow(db.q.SQL(), db.q.Args()...)
}

func (db QueryRow) SQL() string {
	return db.q.SQL()
}

func (db QueryRow) Args() []interface{} {
	return db.q.Args()
}
