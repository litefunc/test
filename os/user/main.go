package main

import (
	"VodoPlay/logger"
	"os"
)

func main() {

	logger.Debug(os.UserCacheDir())
	logger.Debug(os.UserConfigDir())
	logger.Debug(os.UserHomeDir())
}
