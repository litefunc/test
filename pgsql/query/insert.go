package query

import (
	"fmt"
	"strings"
)

func (q Query) Insert(tb string, cols []string, args ...interface{}) Query {

	var vals []string
	for i := range cols {
		vals = append(vals, fmt.Sprintf("$%d", i+1))
	}

	sql := fmt.Sprintf(`%s %s(%s) VALUES(%s)`, INSERT.String(), tb, strings.Join(cols, ", "), strings.Join(vals, ", "))
	q.sql = append(q.sql, sql)
	q.args = append(q.args, args...)

	return q
}

func (q Query) Returning(cols ...string) Query {
	sql := fmt.Sprintf(`RETURNING %s`, strings.Join(cols, ", "))
	q.sql = append(q.sql, sql)
	return q
}
