package main

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"os"
	"test/model/cmd/gorm/internal"

	"github.com/jinzhu/gorm"
)

func main() {
	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config
	dbConfig := config.GetPgsqlConfig(cfg.DB)
	logger.Debug(dbConfig)
	// db, err := gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
	db, err := gorm.Open("postgres", dbConfig)
	if err != nil {
		logger.Error(err)
		return
	}
	defer db.Close()

	db.SingularTable(true)

	internal.Portal(db)

	internal.Tx(db)

}
