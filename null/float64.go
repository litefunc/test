package null

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// Float64 is an alias for sql.NullFloat64 data type
type Float64 struct {
	sql.NullFloat64
}

// MarshalJSON for Float64
func (nf Float64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for Float64
func (nf *Float64) UnmarshalJSON(b []byte) error {
	var err error

	if strings.Contains(string(b), "Valid") {
		err = json.Unmarshal(b, &nf.NullFloat64)
	} else {
		err = json.Unmarshal(b, &nf.Float64)
	}

	nf.Valid = (err == nil)
	return err
}

func NewFloat64(i float64) Float64 {
	x := sql.NullFloat64{Float64: i, Valid: true}
	return Float64{x}
}
