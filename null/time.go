package null

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/lib/pq"
)

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
