package null

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/lib/pq"
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
	} else if string(b) == "null"{
		ni.Valid = false
		return nil
	}else {
		err = json.Unmarshal(b, &ni.Int64)
	}

	ni.Valid = (err == nil)
	return err
}

func NewInt64(i int64) Int64 {
	x := sql.NullInt64{Int64: i, Valid: true}
	return Int64{x}
}

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

// Time is an alias for pq.NullTime data type
type Time struct {
	pq.NullTime
}

// MarshalJSON for Time
func (nt Time) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return nt.Time.MarshalJSON()
}

// UnmarshalJSON for Time
func (ns *Time) UnmarshalJSON(b []byte) error {
	var err error

	if strings.Contains(string(b), "Valid") {
		err = json.Unmarshal(b, &ns.NullTime)
	} else {
		err = json.Unmarshal(b, &ns.Time)
	}

	ns.Valid = (err == nil)
	return err
}

func NewTime(t time.Time) Time {
	x := pq.NullTime{Time: t, Valid: true}
	return Time{x}
}
