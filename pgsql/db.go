package pgsql

import (
	"cloud/lib/logger"
	"context"
	"database/sql"
	"fmt"
	"strings"
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

	return NewDB(db), nil
}

func NewDB(db *sqlx.DB) *DB {
	db.MapperFunc(toSnakeCase)
	return &DB{DB: db}
}

func (db DB) SqlxDB() *sqlx.DB {
	return db.DB
}

func (db DB) SqlDB() *sql.DB {
	return db.DB.DB
}

func (db DB) Close() error {
	if err := db.DB.Close(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (db DB) Begin() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return NewTx(tx), nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return NewTx(tx), nil
}

func (db DB) Insert(md interface{}) Exec {

	var cols []string
	var vals []interface{}

	s := GetSerial(md)

	cols1 := GetCols(md)
	vals1 := GetValues(md)

	for i := range cols1 {
		if cols1[i] != s {
			cols = append(cols, cols1[i])
			vals = append(vals, vals1[i])
		}
	}

	return NewExec(db.DB, db.q.Insert(GetTable(md), cols, vals...), md)
}

func (db DB) Select(md interface{}, cols ...string) Query {

	var q Query
	if len(cols) != 0 {
		q = NewQuery(db.DB, db.q.Select(GetTable(md), cols...), md)
	} else {
		q = NewQuery(db.DB, db.q.Select(GetTable(md), GetCols(md)...), md)
	}

	return q
}

func (db DB) SelectByPk(md interface{}, cols ...string) Query {

	pks := GetPks(md)
	cvs := GetColsValues(md)

	var pcols []string
	var args []interface{}
	for _, k := range pks {
		pcols = append(pcols, fmt.Sprintf(`%s=?`, k))
		args = append(args, cvs[k])
	}

	var q query.Query
	if len(cols) != 0 {
		q = db.q.Select(GetTable(md), cols...).Where(strings.Join(pcols, " AND "), args...)
	} else {
		q = db.q.Select(GetTable(md), GetCols(md)...).Where(strings.Join(pcols, " AND "), args...)
	}

	return NewQuery(db.DB, q, md)
}

func (db DB) Update(md interface{}) Exec {
	q := db.q.Update(GetTable(md), nil)
	return NewExec(db.DB, q, md)
}

func (db DB) UpdateByPk(md interface{}) Exec {
	pks := GetPks(md)
	cvs := GetColsValues(md)

	var cols []string
	var args []interface{}
	for _, k := range pks {
		cols = append(cols, fmt.Sprintf(`%s=?`, k))
		args = append(args, cvs[k])
		delete(cvs, k)
	}

	q := db.q.Update(GetTable(md), cvs).Where(strings.Join(cols, " AND "), args...)

	return NewExec(db.DB, q, md)
}

func (db DB) Delete(md interface{}) Exec {
	return NewExec(db.DB, db.q.Delete(GetTable(md)), md)
}

func (db DB) DeleteByPk(md interface{}) Exec {

	pks := GetPks(md)
	cvs := GetColsValues(md)

	var cols []string
	var args []interface{}
	for _, k := range pks {
		cols = append(cols, fmt.Sprintf(`%s=?`, k))
		args = append(args, cvs[k])
	}

	q := db.q.Delete(GetTable(md)).Where(strings.Join(cols, " AND "), args...)

	return NewExec(db.DB, q, md)
}

func (db DB) Truncate(md interface{}) Exec {
	return NewExec(db.DB, db.q.Truncate(GetTable(md)), md)
}

func (db DB) Count(md interface{}) QueryRow {

	q := NewQueryRow(db.DB, db.q.Count(GetTable(md)))
	return q
}

func (db DB) CountColumn(md interface{}, col string) QueryRow {

	q := NewQueryRow(db.DB, db.q.CountColumn(GetTable(md), col))
	return q
}

func (db DB) CountDistinct(md interface{}, col string) QueryRow {

	q := NewQueryRow(db.DB, db.q.CountDistinct(GetTable(md), col))
	return q
}
