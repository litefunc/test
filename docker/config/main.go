package main

import (
	"VodoPlay/logger"
	"os"
	"os/user"
	"path"
)

func main() {

	p := path.Join(os.Getenv("HOME"), ".docker/config.json")
	logger.Debug(p)

	user, err := user.Current()
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(user)

}
