package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Bool struct {
	v sql.NullBool
}

func NewBool(v bool) Bool {

	return Bool{
		v: sql.NullBool{Bool: v, Valid: true},
	}
}

func (rec Bool) MarshalJSON() ([]byte, error) {
	if !rec.v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(rec.v.Bool)
}

func (rec *Bool) UnmarshalJSON(b []byte) error {
	if string(b) == `null` {
		rec.v.Valid = false
		rec.v.Bool = false
		return nil
	}

	err := json.Unmarshal(b, &rec.v.Bool)

	rec.v.Valid = (err == nil)
	return err
}

func (rec Bool) Valid() bool {
	return rec.v.Valid
}

func (rec Bool) Bool() bool {
	return rec.v.Bool
}

// Scan implements the sql.Scanner interface.
func (rec *Bool) Scan(v interface{}) error {
	return rec.v.Scan(v)
}

// Value implements the sql/driver Valuer interface.
func (rec Bool) Value() (driver.Value, error) {
	return rec.v.Value()
}
