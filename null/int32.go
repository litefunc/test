package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Int32 struct {
	v sql.NullInt32
}

func NewInt32(v int32) Int32 {

	return Int32{
		v: sql.NullInt32{Int32: v, Valid: true},
	}
}

func (rec Int32) MarshalJSON() ([]byte, error) {
	if !rec.v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(rec.v.Int32)
}

func (rec *Int32) UnmarshalJSON(b []byte) error {
	if string(b) == `null` {
		rec.v.Valid = false
		rec.v.Int32 = 0
		return nil
	}

	err := json.Unmarshal(b, &rec.v.Int32)

	rec.v.Valid = (err == nil)
	return err
}

func (rec Int32) Valid() bool {
	return rec.v.Valid
}

func (rec Int32) Int32() int32 {
	return rec.v.Int32
}

// Scan implements the sql.Scanner interface.
func (rec *Int32) Scan(v interface{}) error {
	return rec.v.Scan(v)
}

// Value implements the sql/driver Valuer interface.
func (rec Int32) Value() (driver.Value, error) {
	return rec.v.Value()
}
