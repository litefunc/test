package null

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
)

func TestInt64Marshall(t *testing.T) {

	var tests = []struct {
		name string
		in   Int64
		want string
	}{
		{"a", Int64{}, `null`},
		{"b", NewInt64(0), `0`},
		{"c", NewInt64(1), `1`},
		{"d", Int64{v: sql.NullInt64{Int64: 1, Valid: false}}, `null`},
		{"e", Int64{v: sql.NullInt64{Int64: 0, Valid: true}}, `0`},
		{"f", Int64{v: sql.NullInt64{Int64: 1, Valid: true}}, `1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsequal(t, tt.in, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSuccessInt64Unmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
		want Int64
	}{
		{"a", `null`, Int64{}},
		{"b", `0`, NewInt64(0)},
		{"c", `1`, NewInt64(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Int64
			um(t, []byte(tt.in), &o)
			if !reflect.DeepEqual(tt.want, o) {
				t.Errorf(`want:%+v, got:%+v`, tt.want, o)
			}
		})
	}
}

func TestFailInt64Unmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
	}{
		{"string", `"s"`},
		{"bool", `true`},
		{"time", `"0001-01-01T00:00:00Z"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Int64
			if err := json.Unmarshal([]byte(tt.in), &o); err == nil {
				t.Error("err is nil")
			}
		})
	}
}
