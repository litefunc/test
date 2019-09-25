package main

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"cloud/server/ota/config"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type MsFws []MsFw

type MsFw struct {
	TableName struct{} `json:"-" sql:"cloud.msfw"`
	ID        uint64   `json:"id" sql:",pk"`
	MsFwUnique
	Bucket string      `json:"bucket"`
	Obj    string      `json:"obj"`
	Time   time.Time   `json:"time"`
	Tag    null.String `json:"tag"`
}

type MsFwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver" sql:",notnull"`
}

type Account struct {
	Bio
	Pwd string `json:"pwd"`
}

type Bio struct {
	ID      uint64      `json:"id"`
	Email   string      `json:"email"`
	Name    string      `json:"name"`
	Company uint64      `json:"com_id" db:"com"`
	Admin   bool        `json:"admin"`
	Type    string      `json:"type"`
	Phone   string      `json:"phone"`
	Note    null.String `json:"note"`
	Photo   null.String `json:"photo"`
}

type Accounts []Account

func main() {
	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config
	dbConfig := config.GetPgsqlConfig(cfg.DB)
	db, err := sqlx.Connect("postgres", dbConfig)
	if err != nil {
		logger.Error(err)
		return
	}
	var mds Accounts
	if err := db.Select(&mds, "SELECT * FROM cloud.account"); err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(mds)
	mds = Accounts{}
	if err := db.Select(&mds, "SELECT com, type FROM cloud.account"); err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(mds)
}
