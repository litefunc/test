package internal

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"database/sql"
	"os"
)

func open(dbConfig string, n int) {
	for i := 0; i < n; i++ {
		db, err := sql.Open("postgres", dbConfig)
		if err != nil {
			logger.Error("db connection", err)
			return
		}
		logger.Debug(i, db)
		// defer db.Close()

		if err := db.Ping(); err != nil {
			logger.Error(err)
		}
		logger.Debug(i, db)
	}
}

func begin(db *sql.DB, n int) {
	for i := 0; i < n; i++ {
		_, err := db.Begin()
		if err != nil {
			logger.Error("db connection", err)
			return
		}
	}
}

func maxIdleConns(db *sql.DB, idle, n int) {
	if idle > -1 {
		db.SetMaxIdleConns(idle)
	}

	for i := 0; i < n; i++ {
		tx, err := db.Begin()
		if err != nil {
			logger.Error("db connection", err)
			return
		}
		defer tx.Rollback()
	}
}

func copy(db *sql.DB, n int) {
	logger.Debugf(`%p`, db)
	for i := 0; i < n; i++ {
		db1 := &*db
		logger.Debugf(`%p`, db1)
		if err := db1.Ping(); err != nil {
			logger.Error(err)
		}
	}
}

func copyv(db *sql.DB, n int) {
	logger.Debugf(`%p`, db)
	for i := 0; i < n; i++ {
		db1 := *db
		if err := db1.Ping(); err != nil {
			logger.Error(err)
		}
		_, err := db1.Begin()
		if err != nil {
			logger.Error("db connection", err)
			return
		}
	}
}

func dbConfig() string {
	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config

	dbConfig := config.GetPgsqlConfig(cfg.DB)
	return dbConfig
}

func Conn() {

	dbConfig := dbConfig()

	// conn(dbConfig, 5)

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		logger.Error("db connection", err)
		return
	}
	defer db.Close()

	begin(db, 5)
	// maxIdleConns(db, 5, 10)
	// copyv(db, 10)

	// read(db, 0)
	// read(db1, 1)
	// read(&db2, 2)
	// read(&db3, 3)

	var wc chan int
	<-wc
}
