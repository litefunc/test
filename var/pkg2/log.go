package pkg2

import (
	"test/logger"
	"test/var/var1"
)

func Log(i int) {
	logger.Debug(2, i, var1.T)
}
