package main

import (
	"cloud/lib/logger"
	"cloud/lib/null"

	"test/pgsql"
	"time"

	"cloud/server/ota/config"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type MsFws []MsFw

type MsFw struct {
	TableName  struct{} `json:"-" sql:"cloud.msfw"`
	ID         uint64   `json:"id" sql:",pk"`
	MsFwUnique `sql:"embed"`
	Bucket     string      `json:"bucket"`
	Obj        string      `json:"obj"`
	Time       time.Time   `json:"time"`
	Tag        null.String `json:"tag"`
}

type MsFwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver" sql:",notnull"`
}

func NewMsFw(version string, com uint64, bk, obj string, time time.Time, tag string) MsFw {

	u := MsFwUnique{Version: version, Com: com}
	return MsFw{
		MsFwUnique: u,
		Bucket:     bk,
		Obj:        obj,
		Time:       time,
		Tag:        null.NewString(tag),
	}
}

type Account struct {
	Bio `sql:"embed"`
	Pwd string `json:"pwd"`
}

type Bio struct {
	ID      uint64      `json:"id"`
	Email   string      `json:"email"`
	Name    string      `json:"name"`
	Company uint64      `json:"com_id" sql:"com" db:"com"`
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
	logger.Debug(dbConfig)
	db, _ := pgsql.Connect(dbConfig)
	defer db.Close()

	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	mds := MsFws{md1, md2}


	db.Select(md1).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(&md1).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(mds).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(&mds).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(nil).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()

	mds1 := Accounts{}

	if err := db.Select(&mds1).Run();err !=nil {
		logger.Error(err)
	}
	for i := range mds1 {
		logger.Debug(mds1[i])
	}

	mds1 = Accounts{}
	if err := db.DB.Select(&mds1, "SELECT * FROM account");err !=nil {
		logger.Error(err)
	}
	for i := range mds1 {
		logger.Debug(mds1[i])
	}
	
}
