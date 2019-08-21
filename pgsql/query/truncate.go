package query

import (
	"fmt"
)

func (q Query) Truncate(tb string) Query {

	sql := fmt.Sprintf(`%s %s`, TRUNCATE.String(), tb)
	q.sql = append(q.sql, sql)

	return q
}
