package internal

import (
	"test/logger"
	"time"
)

func W() {
	now := time.Now().Add(time.Hour * 24 * 365)
	logger.Debug(now.UTC())
}
