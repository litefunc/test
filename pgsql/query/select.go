package query

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func (q Query) Select(tb string, cols ...string) Query {
	q.qt = SELECT
	q.sql = append(q.sql, SELECT.String())
	q.sql = append(q.sql, strings.Join(cols, ", "))
	q.sql = append(q.sql, "FROM", tb)
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
