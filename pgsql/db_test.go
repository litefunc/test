package pgsql

import (
	"cloud/lib/logger"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
)

func TestBasicCRUD(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")

	db, _ := Connect(dbConfig)
	defer db.Close()

	var mds []Tb
	db.Select(&mds).SQL()
	if err := db.Select(&mds).Run(); err != nil {
		t.Error(err)
		return
	}

	// len(mds) should be 0
	if err := ModelsEqual(db, mds); err != nil {
		t.Error(err)
		return
	}

	if err := db.Insert(md1).Returning("id").Run().Scan(&md1.A); err != nil {
		t.Error(err)
		return
	}
	if err := db.Insert(md2).Returning("id").Run().Scan(&md2.A); err != nil {
		t.Error(err)
		return
	}

	md := Tb{A: md1.A}
	if err := db.SelectByPk(&md).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(md1, md); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(db, append(mds, md1, md2)); err != nil {
		t.Error(err)
		return
	}

	if _, err := db.DeleteByPk(md2).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelsEqual(db, append(mds, md1)); err != nil {
		t.Error(err)
		return
	}

	if _, err := db.Truncate(mds).Run(); err != nil {
		t.Error(err)
		return
	}

	db.Select(&mds).SQL()
	if err := db.Select(&mds).Run(); err != nil {
		t.Error(err)
		return
	}
	// len(mds) should be 0
	if err := ModelsEqual(db, mds); err != nil {
		t.Error(err)
		return
	}

}

func ModelsEqual(db *DB, want Tbs) error {
	var got Tbs
	if err := db.Select(&got).Run(); err != nil {
		logger.Error(err)
		return err
	}

	if len(want) != len(got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}

func ModelEqual(want, got Tb) error {

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}
