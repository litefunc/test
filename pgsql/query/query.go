package query

import (
	"cloud/lib/logger"
	"fmt"
	"regexp"
	"strings"
)

type Query struct {
	qt   QueryType
	sql  []string
	args []interface{}
}

func New() Query {
	return Query{}
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

func (q Query) Type() QueryType {
	return q.qt
}

type QueryType int

const (
	NONE QueryType = iota
	SELECT
	INSERT
	UPDATE
	DELETE
	COUNT
	DISTINCT
	TRUNCATE
)

func (qt QueryType) String() string {
	switch qt {
	case SELECT:
		return "SELECT"
	case INSERT:
		return "INSERT INTO"
	case UPDATE:
		return "UPDATE"
	case DELETE:
		return "DELETE"
	case COUNT:
		return "COUNT"
	case DISTINCT:
		return "SELECT DISTINCT"
	case TRUNCATE:
		return "TRUNCATE TABLE"
	default:
		return ""
	}
}
