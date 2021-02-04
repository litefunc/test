package null

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// Bool is an alias for sql.NullBool data type
type Bool struct {
	sql.NullBool
}

// MarshalJSON for Bool
func (nb Bool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON for Bool
func (nb *Bool) UnmarshalJSON(b []byte) error {

	var err error

	if strings.Contains(string(b), "Valid") {
		err = json.Unmarshal(b, &nb.NullBool)
	} else {
		err = json.Unmarshal(b, &nb.Bool)
	}

	nb.Valid = (err == nil)
	return err
}

func NewBool(b bool) Bool {
	x := sql.NullBool{Bool: b, Valid: true}
	return Bool{x}
}
