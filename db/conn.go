package main

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"database/sql"
	"fmt"
	"os"
)

func read(db *sql.DB, i int) {
	rows, err := db.Query(`SELECT * FROM cloud.device_stat_vod`)
	if err != nil {
		logger.Error(err)
		return
	}
	defer rows.Close()
	fmt.Println(i)
}

func main() {

	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config

	dbConfig := config.GetPgsqlConfig(cfg.DB)

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		logger.Error("db connection", err)
		return
	}
	defer db.Close()

	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		logger.Error(err)
	}
	if err := db.Ping(); err != nil {
		logger.Error(err)
	}
	db1 := db
	if err := db1.Ping(); err != nil {
		logger.Error(err)
	}

	db2 := *db
	if err := db2.Ping(); err != nil {
		logger.Error(err)
	}

	db3 := *db
	if err := db3.Ping(); err != nil {
		logger.Error(err)
	}

	tx, err := db.Begin()
	if err != nil {
		logger.Error(err)
	}
	defer tx.Commit()

	tx1, err := db1.Begin()
	if err != nil {
		logger.Error(err)
	}
	defer tx1.Commit()

	tx2, err := db2.Begin()
	if err != nil {
		logger.Error(err)
	}
	defer tx2.Commit()

	tx3, err := db3.Begin()
	if err != nil {
		logger.Error(err)
	}
	defer tx3.Commit()

	tx.Commit()
	tx1.Commit()
	tx2.Commit()
	tx3.Commit()

	read(db, 0)
	read(db1, 1)
	read(&db2, 2)
	read(&db3, 3)

	var wc chan int
	<-wc
}
