package pkg1

import (
	"test/logger"
	"test/var/var1"
)

var d = logger.DefaultLogger()

var l = logger.NewLogger(d.Flag, d.Level, true)

func Log(i int) {
	logger.Trace(1, i, var1.T)
	l.Info(1, i, var1.T)
}
