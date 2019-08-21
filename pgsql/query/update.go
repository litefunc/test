package query

import (
	"fmt"
	"strings"
)

func (q Query) Update(tb string, set map[string]interface{}) Query {
	if set != nil {
		var cols []string
		var args []interface{}
		for k := range set {
			cols = append(cols, fmt.Sprintf(`%s=?`, k))
			args = append(args, set[k])
		}

		sql := fmt.Sprintf(`%s %s SET %s`, UPDATE.String(), tb, strings.Join(cols, ", "))
		q.sql = append(q.sql, sql)
		q.args = append(q.args, args...)
		return q
	}

	sql := fmt.Sprintf(`%s %s`, UPDATE.String(), tb)
	q.sql = append(q.sql, sql)
	return q
}

func (q Query) Set(condition string, args ...interface{}) Query {
	q.sql = append(q.sql, "SET", condition)
	q.args = append(q.args, args...)
	return q
}
