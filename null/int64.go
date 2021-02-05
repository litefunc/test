package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Int64 struct {
	v sql.NullInt64
}

func NewInt64(v int64) Int64 {

	return Int64{
		v: sql.NullInt64{Int64: v, Valid: true},
	}
}

func (rec Int64) MarshalJSON() ([]byte, error) {
	if !rec.v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(rec.v.Int64)
}

func (rec *Int64) UnmarshalJSON(b []byte) error {
	if string(b) == `null` {
		rec.v.Valid = false
		rec.v.Int64 = 0
		return nil
	}

	err := json.Unmarshal(b, &rec.v.Int64)

	rec.v.Valid = (err == nil)
	return err
}

func (rec Int64) Valid() bool {
	return rec.v.Valid
}

func (rec Int64) Int64() int64 {
	return rec.v.Int64
}

// Scan implements the sql.Scanner interface.
func (rec *Int64) Scan(v interface{}) error {
	return rec.v.Scan(v)
}

// Value implements the sql/driver Valuer interface.
func (rec Int64) Value() (driver.Value, error) {
	return rec.v.Value()
}
