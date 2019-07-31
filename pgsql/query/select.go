package query

import (
	"cloud/lib/logger"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
)

type Query struct {
	sql  []string
	args []interface{}
}

func New() Query {
	return Query{}
}

func (q Query) Select(cols ...string) Query {
	q.sql = append(q.sql, "SELECT")
	q.sql = append(q.sql, strings.Join(cols, ", "))
	return q
}

func (q Query) From(table string) Query {
	q.sql = append(q.sql, "FROM")
	q.sql = append(q.sql, table)
	return q
}

func (q Query) Where(condition string, args ...interface{}) Query {

	q.sql = append(q.sql, "WHERE", condition)
	q.args = append(q.args, args...)
	return q
}

func (q Query) Order(args ...string) Query {

	q.sql = append(q.sql, "ORDER BY")
	q.sql = append(q.sql, strings.Join(args, ", "))
	return q
}

func (q Query) Limit(n uint64) Query {

	q.sql = append(q.sql, fmt.Sprintf("LIMIT %d", n))
	return q
}

func (q Query) SQL() string {

	sql := strings.Join(q.sql, " ")
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
	logger.Debug(sql, q.args)
	return sql
}

func (q Query) Args() []interface{} {
	return q.args
}
