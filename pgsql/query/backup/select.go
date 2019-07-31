package query

import (
	"cloud/lib/logger"
	"fmt"
	"strings"
)

type SQL struct {
	Sql  string
	cols []string
}

func (self SQL) String() string {
	return self.Sql
}

func (self SQL) Log() {
	logger.Debug(self.Sql)
}

func (self SQL) append(sql string) SQL {

	self.Sql = self.Sql + " " + sql
	return self
}

func (self SQL) appendCols(cols []string) SQL {

	self.cols = append(self.cols, cols...)
	return self
}

func (self SQL) From(table string) SQL {
	return self.append("FROM").append(table)
}

func (self SQL) Where(str string) SQL {
	if str == "" {
		return self.append("WHERE")
	}
	return self.append("WHERE").append(str)
}

func (self SQL) Columns(columns ...string) SQL {

	n := len(self.cols)
	self = self.appendCols(columns)

	var strs []string
	for i := range columns {
		strs = append(strs, fmt.Sprintf("%s=$%d", columns[i], n+i+1))
	}
	q := strings.Join(strs, " AND ")
	return self.append(q)
}

func (self SQL) WhereColumns(columns ...string) SQL {

	n := len(self.cols)
	self = self.appendCols(columns)

	var strs []string
	for i := range columns {
		strs = append(strs, fmt.Sprintf("%s=$%d", columns[i], n+i+1))
	}
	q := strings.Join(strs, " AND ")
	return self.append("WHERE").append(q)
}

func (self SQL) And(str string) SQL {
	if str == "" {
		return self.append("AND")
	}
	return self.append("AND").append(str)
}

func (self SQL) AndColumns(columns ...string) SQL {

	n := len(self.cols)
	self = self.appendCols(columns)

	var strs []string
	for i := range columns {
		strs = append(strs, fmt.Sprintf("%s=$%d", columns[i], n+i+1))
	}
	q := strings.Join(strs, " AND ")
	return self.append("AND").append(q)
}

func (self SQL) Limit(i int) SQL {
	return self.append("LIMIT").append(fmt.Sprintf(`%d`, i))
}

func Select(columns ...string) SQL {

	sql := SQL{Sql: "SELECT"}
	q := strings.Join(columns, ", ")

	return sql.append(q)
}

func Distinct(columns ...string) SQL {

	sql := SQL{Sql: "SELECT DISTINCT"}
	q := strings.Join(columns, ", ")

	return sql.append(q)
}

func CountAll(table string) SQL {

	sql := SQL{Sql: "SELECT count(*)"}.From(table)
	return sql
}
