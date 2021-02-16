package main

import "test/logger"

func main() {
	var i int
	logger.Debugf(`%p`, &i)

	var ip *int
	logger.Debugf(`%p`, ip)
	ip = &i
	logger.Debug(ip)
	logger.Debugf(`%p`, ip)
	logger.Debugf(`%p`, &ip)
}
