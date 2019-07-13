package internal

import (
	"cloud/lib/logger"
	"context"
	"time"

	"github.com/go-xorm/xorm"
)

func Tx(engine *xorm.Engine) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	session := engine.NewSession()
	go func() {
		time.Sleep(time.Second * 10)
		logger.Error(session.Commit())
		logger.Error(session.Rollback())
	}()

	// add Begin() before any action
	logger.Error(session.Begin())
	logger.Error(session.Begin())

	session1 := engine.Context(ctx)
	defer session1.Close()

	logger.Error(session1.Begin())

	var wc chan struct{}
	<-wc
}
