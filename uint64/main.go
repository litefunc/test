package main

import "test/logger"

func main() {

	var i uint64 = 1
	var j uint64 = 2
	logger.Debug(i - j)

}
