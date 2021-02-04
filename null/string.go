package null

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// String is an alias for sql.NullString data type
type String struct {
	sql.NullString
}

// MarshalJSON for String
func (ns String) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for String
func (ns *String) UnmarshalJSON(b []byte) error {
	var err error

	if strings.Contains(string(b), "Valid") {
		err = json.Unmarshal(b, &ns.NullString)
	} else {
		err = json.Unmarshal(b, &ns.String)
	}

	ns.Valid = (err == nil)
	return err
}

func NewString(s string) String {
	x := sql.NullString{String: s, Valid: true}
	return String{x}
}
