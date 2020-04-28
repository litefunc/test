package internal

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"cloud/server/ota/model/v1"

	"context"
	"os"
	"time"

	"github.com/go-pg/pg"
)

func pgRead(tx *pg.Tx, i int) {
	var md model.VideoShareWithCompany
	if err := tx.Model(&md).Select(); err != nil {
		logger.Error(err)
	}
	logger.Debug(i)
}

func PgRead() {
	p := os.Getenv("GOPATH") + "/src/cloud/server/ota/config/config.local.json"
	config.ParseConfig(p, &config.Config)
	cfg := &config.Config

	// dbConfig := config.GetPgsqlConfig(cfg.DB)
	// db, err := database.Connect(dbConfig)
	// if err != nil {
	// 	return
	// }
	// defer db.Close()

	// logger.SetLogger(cfg.Logger.Flag, cfg.Logger.Level, db, cfg.Logger.SaveToDir, cfg.Logger.Service)

	pgdb, err := cfg.DB.Pg().Connect()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(pgdb)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	conn := pgdb.Conn().WithContext(ctx)
	tx, err := conn.Begin()
	if err != nil {
		logger.Error(err)
		return
	}

	cancel()
	for i := 0; i < 10; i++ {
		pgRead(tx, i)
		time.Sleep(time.Second)
	}
	logger.Error(tx.Commit())

	var wc chan int
	<-wc

}
