package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type String struct {
	v sql.NullString
}

func NewString(v string) String {

	return String{
		v: sql.NullString{String: v, Valid: true},
	}
}

func (rec String) MarshalJSON() ([]byte, error) {
	if !rec.v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(rec.v.String)
}

func (rec *String) UnmarshalJSON(b []byte) error {
	if string(b) == `null` {
		rec.v.Valid = false
		rec.v.String = ``
		return nil
	}

	err := json.Unmarshal(b, &rec.v.String)

	rec.v.Valid = (err == nil)
	return err
}

func (rec String) Valid() bool {
	return rec.v.Valid
}

func (rec String) String() string {
	return rec.v.String
}

// Scan implements the sql.Scanner interface.
func (rec *String) Scan(v interface{}) error {
	return rec.v.Scan(v)
}

// Value implements the sql/driver Valuer interface.
func (rec String) Value() (driver.Value, error) {
	return rec.v.Value()
}
