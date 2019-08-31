package query

import (
	"fmt"
)

func (q Query) Delete(tb string) Query {
	q.qt = DELETE
	sql := fmt.Sprintf(`%s FROM %s`, DELETE.String(), tb)
	q.sql = append(q.sql, sql)

	return q
}
