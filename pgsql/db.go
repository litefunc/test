package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

type DB struct {
	db   *sql.DB
	sql  []string
	args []interface{}
}

func (db DB) Select(cols ...string) DB {
	db.sql = append(db.sql, "SELECT")
	db.sql = append(db.sql, strings.Join(cols, ", "))
	return db
}

func (db DB) From(table string) DB {
	db.sql = append(db.sql, "FROM")
	db.sql = append(db.sql, table)
	return db
}

func (db DB) Where(condition string, args ...interface{}) DB {

	db.sql = append(db.sql, "WHERE", condition)
	db.args = append(db.args, args...)
	return db
}

func (db DB) Order(args ...string) DB {

	db.sql = append(db.sql, "ORDER BY")
	db.sql = append(db.sql, strings.Join(args, ", "))
	return db
}

func (db DB) Limit(n uint64) DB {

	db.sql = append(db.sql, fmt.Sprintf("LIMIT %d", n))
	return db
}

func (db DB) SQL() string {

	sql := strings.Join(db.sql, " ")
	sqls := strings.Split(sql, " ")

	var sqls1 []string
	for _, str := range sqls {
		if s := strings.Replace(str, " ", "", -1); s != "" {
			sqls1 = append(sqls1, s)
		}
	}

	sql = strings.Join(sqls1, " ")
	sql = strings.Replace(sql, "( ", "(", -1)
	sql = strings.Replace(sql, " )", ")", -1)
	sql = strings.Replace(sql, " in (?)", " = ANY(?)", -1)
	sql = strings.Replace(sql, " in(?)", " = ANY(?)", -1)
	sql = strings.Replace(sql, " IN (?)", " = ANY(?)", -1)
	sql = strings.Replace(sql, " IN(?)", " = ANY(?)", -1)
	r, err := regexp.Compile("\\?")
	if err != nil {
		logger.Error(err)
	}

	qs := r.FindAllString(sql, -1)
	for i := range qs {
		sql = strings.Replace(sql, "?", fmt.Sprintf("$%d", i+1), 1)
	}
	logger.Debug(sql, db.args)
	return sql
}

func (db DB) Args() []interface{} {
	return db.args
}
