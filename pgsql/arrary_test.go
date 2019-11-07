package pgsql

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lib/pq"
)

func TestArray(t *testing.T) {

	type Tba struct {
		TableName struct{}       `json:"-" pg:"test.Tb_array"`
		I         pq.Int64Array  `json:"i" pg:"i"`
		T         pq.StringArray `json:"t" pg:"t"`
	}

	type Tbas []Tba

	newTba := func(i []int64, t []string) Tba {
		return Tba{I: i, T: t}
	}

	type Tbb struct {
		TableName struct{} `json:"-" pg:"test.Tb_array"`
		I         []int    `json:"i" pg:"i"`
		T         []string `json:"t" pg:"t"`
	}

	type Tbbs []Tbb

	newTbb := func(i []int, t []string) Tbb {
		return Tbb{I: i, T: t}
	}

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()
	test := NewTester(t)

	var mdas Tbas
	test.Run(tx.Select(&mdas))

	is := []int64{1, 2}
	ts := []string{"a", "b"}

	mda1 := newTba(is, ts)
	test.Run(tx.Insert(mda1))
	test.Run(tx.Select(&mdas))

	var mdbs Tbbs
	// test.Run(tx.Select(&mdbs))

	mdb1 := newTbb([]int{1, 2}, ts)
	if err := tx.Insert(mdb1).Run(); err.Error() != errors.New(`sql: converting argument $1 type: unsupported type []int, a slice of int`).Error() {
		t.Error(err)
	}
	if err := tx.Select(&mdbs).Run(); err.Error() != errors.New(`sql: Scan error on column index 0, name "i": unsupported Scan, storing driver.Value type []uint8 into type *[]int`).Error() {
		t.Error(err)
	}
}
