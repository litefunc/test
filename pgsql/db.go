package pgsql

import (
	"cloud/lib/logger"
	"test/pgsql/query"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
	q  query.Query
	md interface{}
}

func Connect(dbConfig string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dbConfig)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &DB{DB: db}, nil
}

func (db DB) Close() error {
	if err := db.DB.Close(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (db DB) Run() error {
	logger.Debug(db.md)
	logger.Debug(db.q.SQL())
	logger.Debug(db.q.Args()...)
	if err := db.DB.Select(db.md, db.q.SQL(), db.q.Args()...); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (db DB) Select(md interface{}) DB {
	db.q = db.q.Select(GetCols(md)...).From(GetTable(md))
	db.md = md
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
