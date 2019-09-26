package pgsql

import (
	"cloud/lib/logger"
	"fmt"
	"strings"
	"test/pgsql/query"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Tx struct {
	*sqlx.Tx
	q  query.Query
	md interface{}
}

func NewTx(tx *sqlx.Tx) *Tx {
	return &Tx{Tx: tx}
}

func (tx Tx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (tx Tx) Rollback() error {
	if err := tx.Tx.Rollback(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (tx Tx) Insert(md interface{}) Exec {

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

	return NewExec(tx.Tx, tx.q.Insert(GetTable(md), cols, vals...), md)
}

func (tx Tx) Select(md interface{}, cols ...string) Query {

	var q Query
	if len(cols) != 0 {
		q = NewQuery(tx.Tx, tx.q.Select(GetTable(md), cols...), md)
	} else {
		q = NewQuery(tx.Tx, tx.q.Select(GetTable(md), GetCols(md)...), md)
	}

	return q
}

func (tx Tx) SelectByPk(md interface{}, cols ...string) Query {

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
		q = tx.q.Select(GetTable(md), cols...).Where(strings.Join(pcols, " AND "), args...)
	} else {
		q = tx.q.Select(GetTable(md), GetCols(md)...).Where(strings.Join(pcols, " AND "), args...)
	}

	return NewQuery(tx.Tx, q, md)
}

func (tx Tx) Update(md interface{}) Exec {
	q := tx.q.Update(GetTable(md), nil)
	return NewExec(tx.Tx, q, md)
}

func (tx Tx) UpdateByPk(md interface{}) Exec {
	pks := GetPks(md)
	cvs := GetColsValues(md)

	var cols []string
	var args []interface{}
	for _, k := range pks {
		cols = append(cols, fmt.Sprintf(`%s=?`, k))
		args = append(args, cvs[k])
		delete(cvs, k)
	}

	q := tx.q.Update(GetTable(md), cvs).Where(strings.Join(cols, " AND "), args...)

	return NewExec(tx.Tx, q, md)
}

func (tx Tx) Delete(md interface{}) Exec {
	return NewExec(tx.Tx, tx.q.Delete(GetTable(md)), md)
}

func (tx Tx) DeleteByPk(md interface{}) Exec {

	pks := GetPks(md)
	cvs := GetColsValues(md)

	var cols []string
	var args []interface{}
	for _, k := range pks {
		cols = append(cols, fmt.Sprintf(`%s=?`, k))
		args = append(args, cvs[k])
	}

	q := tx.q.Delete(GetTable(md)).Where(strings.Join(cols, " AND "), args...)

	return NewExec(tx.Tx, q, md)
}

func (tx Tx) Truncate(md interface{}) Exec {
	return NewExec(tx.Tx, tx.q.Truncate(GetTable(md)), md)
}

func (tx Tx) Count(md interface{}) QueryRow {

	q := NewQueryRow(tx.Tx, tx.q.Count(GetTable(md)))
	return q
}

func (tx Tx) CountColumn(md interface{}, col string) QueryRow {

	q := NewQueryRow(tx.Tx, tx.q.CountColumn(GetTable(md), col))
	return q
}

func (tx Tx) CountDistinct(md interface{}, col string) QueryRow {

	q := NewQueryRow(tx.Tx, tx.q.CountDistinct(GetTable(md), col))
	return q
}
