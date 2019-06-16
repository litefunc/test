package main

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

func read(tx *sql.Tx, i int) {
	rows, err := tx.Query(`SELECT * FROM cloud.device_stat_vod`)
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
	// cfg.DB.Password = "fake"

	dbConfig := config.GetPgsqlConfig(cfg.DB)

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		logger.Error("db connection", err)
		return
	}
	defer db.Close()

	ctx, _ := context.WithTimeout(context.Background(), 4*time.Second)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err)
		logger.Debug(tx == nil)
		logger.Error(tx.Rollback())
	}
	logger.Debug(tx)

	for i := 0; i < 10; i++ {
		read(tx, i)
		time.Sleep(time.Second)
	}

	logger.Debug(tx)

	var wc chan int
	<-wc
}
