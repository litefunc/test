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

func (q Query) SelectDistinct(tb string, cols ...string) Query {
	q.qt = DISTINCT
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

func (q Query) In(arg interface{}) Query {
	q.sql = append(q.sql, "IN (?) ")
	q.args = append(q.args, arg)
	return q
}

func (q Query) GroupBy(cols ...string) Query {
	q.sql = append(q.sql, "GROUP BY", strings.Join(cols, ", "))
	return q
}

func (q Query) Having(condition string, args ...interface{}) Query {
	q.sql = append(q.sql, "HAVING", condition)
	q.args = append(q.args, args...)
	return q
}

func (q Query) OrderBy(args ...string) Query {
	q.sql = append(q.sql, "ORDER BY")
	q.sql = append(q.sql, strings.Join(args, ", "))
	return q
}

func (q Query) Limit(n uint64) Query {
	q.sql = append(q.sql, fmt.Sprintf("LIMIT %d", n))
	return q
}

func (q Query) Offset(n uint64) Query {
	q.sql = append(q.sql, fmt.Sprintf("OFFSET %d", n))
	return q
}

func (q Query) Count(tb string) Query {
	q.qt = SELECT
	q.sql = append(q.sql, SELECT.String())
	q.sql = append(q.sql, "count(*) FROM", tb)
	return q
}

func (q Query) CountColumn(tb string, col string) Query {
	q.qt = SELECT
	q.sql = append(q.sql, SELECT.String())
	q.sql = append(q.sql, fmt.Sprintf("count(%s) FROM %s", col, tb))
	return q
}

func (q Query) CountDistinct(tb string, col string) Query {
	q.qt = SELECT
	q.sql = append(q.sql, SELECT.String())
	q.sql = append(q.sql, fmt.Sprintf("count(DISTINCT %s) FROM %s", col, tb))
	return q
}
