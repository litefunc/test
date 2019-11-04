package model

import (
	"fmt"
	"strings"
	"test/pgsql"
)

type DB struct {
	*pgsql.DB
}

type Tx struct {
	*pgsql.Tx
}

func arg(m map[string]interface{}) (string, []interface{}) {
	var keys []string
	var vals []interface{}
	for k, v := range m {
		keys = append(keys, fmt.Sprintf(`%s=?`, k))
		vals = append(vals, v)
	}
	return strings.Join(keys, " AND "), vals
}
