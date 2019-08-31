package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"test/pgsql/query"
)

type QueryRow struct {
	DB Database
	q  query.Query
}

func NewQueryRow(db Database, q query.Query) QueryRow {
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

func (db QueryRow) Where(condition string, args ...interface{}) QueryRow {
	db.q = db.q.Where(condition, args...)
	return db
}

func (db QueryRow) GroupBy(cols ...string) QueryRow {
	db.q = db.q.GroupBy(cols...)
	return db
}

func (db QueryRow) Having(condition string, args ...interface{}) QueryRow {
	db.q = db.q.Having(condition, args...)
	return db
}

func (db QueryRow) OrderBy(args ...string) QueryRow {
	db.q = db.q.OrderBy(args...)
	return db
}

func (db QueryRow) Limit(n uint64) QueryRow {
	db.q = db.q.Limit(n)
	return db
}

func (db QueryRow) Offset(n uint64) QueryRow {
	db.q = db.q.Offset(n)
	return db
}
