package null

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestTimeMarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   Time
		want string
	}{
		{"a", Time{}, `null`},
		{"b", NewTime(time.Time{}), `"0001-01-01T00:00:00Z"`},
		{"c", NewTime(time.Date(2021, time.February, 5, 0, 0, 0, 0, time.UTC)), `"2021-02-05T00:00:00Z"`},
		{"d", Time{v: sql.NullTime{Time: time.Date(2021, time.February, 5, 0, 0, 0, 0, time.UTC), Valid: false}}, `null`},
		{"e", Time{v: sql.NullTime{Time: time.Date(2021, time.February, 5, 0, 0, 0, 0, time.UTC), Valid: true}}, `"2021-02-05T00:00:00Z"`},
		{"f", Time{v: sql.NullTime{Time: time.Time{}, Valid: true}}, `"0001-01-01T00:00:00Z"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsequal(t, tt.in, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSuccessTimeUnmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
		want Time
	}{
		{"a", `null`, Time{}},
		{"b", `"0001-01-01T00:00:00Z"`, NewTime(time.Time{})},
		{"c", `"2021-02-05T00:00:00Z"`, NewTime(time.Date(2021, time.February, 5, 0, 0, 0, 0, time.UTC))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Time
			um(t, []byte(tt.in), &o)
			if !reflect.DeepEqual(tt.want, o) {
				t.Errorf(`want:%+v, got:%+v`, tt.want, o)
			}
		})
	}
}

func TestFailTimeUnmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
	}{

		{"unquoted", `0001-01-01T00:00:00Z`},
		{"bool", `true`},
		{"int", `1`},
		{"string", `"s"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Time
			if err := json.Unmarshal([]byte(tt.in), &o); err == nil {
				t.Error("err is nil")
			}
		})
	}
}
