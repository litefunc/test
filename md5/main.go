package main

import (
	"cloud/lib/logger"
	"crypto/md5"
)

func main() {
	data := []byte("These pretzels are making me thirsty.")
	logger.Debugf("%x", md5.Sum(data))
}
