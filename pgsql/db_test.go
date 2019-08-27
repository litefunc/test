package pgsql

import (
	"cloud/lib/logger"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func TestBasicCRUD(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")

	db, _ := Connect(dbConfig)
	defer db.Close()

	if _, err := db.Insert(md1).Run(); err != nil {
		t.Error(err)
		return
	}

	if _, err := db.Insert(md2).Run(); err != nil {
		t.Error(err)
		return
	}

	// res, err := db.Insert(md1).Returning("id").Run()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// id, err := res.LastInsertId()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// md1.A = id

	// res, err = db.Insert(md2).Returning("id").Run()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// id, err = res.LastInsertId()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// md2.A = id

	// logger.Debug(md1, md2)

	var mds []Tb
	db.Select(&mds).SQL()
	if err := db.Select(&mds).Run(); err != nil {
		t.Error(err)
		return
	}
	logger.Debug(mds)

	if _, err := db.Delete(mds).Run(); err != nil {
		t.Error(err)
		return
	}

}
