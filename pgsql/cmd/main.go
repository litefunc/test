package main

import (
	"test/pgsql"

	"github.com/lib/pq"
)

func main() {

	var db pgsql.DB
	db.Select("id", "com", "ver").From("cloud.msfw").Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()

}
