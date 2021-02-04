package null

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// Int64 is an alias for sql.NullInt64 data type
type Int64 struct {
	sql.NullInt64
}

// MarshalJSON for Int64
func (ni Int64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for Int64
func (ni *Int64) UnmarshalJSON(b []byte) error {
	var err error

	if strings.Contains(string(b), "Valid") {
		err = json.Unmarshal(b, &ni.NullInt64)
	} else if string(b) == "null" {
		ni.Valid = false
		return nil
	} else {
		err = json.Unmarshal(b, &ni.Int64)
	}

	ni.Valid = (err == nil)
	return err
}

func NewInt64(i int64) Int64 {
	x := sql.NullInt64{Int64: i, Valid: true}
	return Int64{x}
}
