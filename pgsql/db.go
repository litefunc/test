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

func (db DB) Insert(md interface{}) Exec {
	return NewExec(db.DB, db.q.Insert(GetTable(md), GetCols(md), GetValues(md)...), md)
}

func (db DB) Select(md interface{}) Query {
	q := NewQuery(db.DB, db.q.Select(GetTable(md), GetCols(md)...), md)
	return q
}

func (db DB) Update(md interface{}) Exec {
	q := db.q.Update(GetTable(md), nil)
	return NewExec(db.DB, q, md)
}

func (db DB) Delete(md interface{}) Exec {
	return NewExec(db.DB, db.q.Delete(GetTable(md)), md)
}

func (db DB) DeleteAll(md interface{}) Exec {
	return NewExec(db.DB, db.q.Delete(GetTable(md)), md)
}
