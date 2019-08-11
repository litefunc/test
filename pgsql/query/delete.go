package query

import (
	"fmt"
)

func (q Query) Delete(tb string) Query {

	sql := fmt.Sprintf(`%s FROM %s`, DELETE.String(), tb)
	q.sql = append(q.sql, sql)

	return q
}
