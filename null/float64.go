package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Float64 struct {
	v sql.NullFloat64
}

func NewFloat64(v float64) Float64 {

	return Float64{
		v: sql.NullFloat64{Float64: v, Valid: true},
	}
}

func (rec Float64) MarshalJSON() ([]byte, error) {
	if !rec.v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(rec.v.Float64)
}

func (rec *Float64) UnmarshalJSON(b []byte) error {
	if string(b) == `null` {
		rec.v.Valid = false
		rec.v.Float64 = 0
		return nil
	}

	err := json.Unmarshal(b, &rec.v.Float64)

	rec.v.Valid = (err == nil)
	return err
}

func (rec Float64) Valid() bool {
	return rec.v.Valid
}

func (rec Float64) Float64() float64 {
	return rec.v.Float64
}

// Scan implements the sql.Scanner interface.
func (rec *Float64) Scan(v interface{}) error {
	return rec.v.Scan(v)
}

// Value implements the sql/driver Valuer interface.
func (rec Float64) Value() (driver.Value, error) {
	return rec.v.Value()
}
