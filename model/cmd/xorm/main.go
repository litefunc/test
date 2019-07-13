package main

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"os"
	"test/model/cmd/xorm/internal"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

func main() {
	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config
	dbConfig := config.GetPgsqlConfig(cfg.DB)

	engine, err := xorm.NewEngine("postgres", dbConfig)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(engine)

	internal.Portal(engine)
	var wc chan struct{}
	<-wc
}
