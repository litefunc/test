package internal

import (
	"cloud/lib/logger"
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

func Tx(db *gorm.DB) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	tx := db.BeginTx(ctx, nil)

	logger.Debug(tx)
	time.Sleep(time.Second * 5)
	tx.Commit()
	tx.Rollback()
	logger.Debug(tx)
	logger.Debug(tx == nil)

	var wc chan struct{}
	<-wc
}
