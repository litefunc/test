package pgsql

import (
	"cloud/lib/logger"
	"fmt"
	"testing"
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
	logger.Debug(mds)

}
