package main

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/msfw"
	"test/pgsql"
)

func main() {

	cols := []string{"a", "b"}
	tb := "test"
	pgsql.Select(cols...).From(tb).Where("").Columns("id", "com").And("c=$3").Limit(2).Log()
	pgsql.Select(cols...).From(tb).WhereColumns("id", "com").AndColumns("gp", "sn").Limit(2).Log()
	pgsql.CountAll(tb).Where("").Columns("id", "com").Columns("gp", "sn").Log()

	var msfws msfw.MsFws
	logger.Debug(msfws)
	logger.Debug(msfws == nil)
	logger.Debug(msfws.Json())
}
