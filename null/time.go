package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Time struct {
	v sql.NullTime
}

func NewTime(v time.Time) Time {

	return Time{
		v: sql.NullTime{Time: v, Valid: true},
	}
}

func (rec Time) MarshalJSON() ([]byte, error) {
	if !rec.v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(rec.v.Time)
}

func (rec *Time) UnmarshalJSON(b []byte) error {
	if string(b) == `null` {
		rec.v.Valid = false
		rec.v.Time = time.Time{}
		return nil
	}

	err := json.Unmarshal(b, &rec.v.Time)

	rec.v.Valid = (err == nil)
	return err
}

func (rec Time) Valid() bool {
	return rec.v.Valid
}

func (rec Time) Time() time.Time {
	return rec.v.Time
}

// Scan implements the sql.Scanner interface.
func (rec *Time) Scan(v interface{}) error {
	return rec.v.Scan(v)
}

// Value implements the sql/driver Valuer interface.
func (rec Time) Value() (driver.Value, error) {
	return rec.v.Value()
}
