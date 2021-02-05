package null

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
)

func TestInt32Marshall(t *testing.T) {

	var tests = []struct {
		name string
		in   Int32
		want string
	}{
		{"a", Int32{}, `null`},
		{"b", NewInt32(0), `0`},
		{"c", NewInt32(1), `1`},
		{"d", Int32{v: sql.NullInt32{Int32: 1, Valid: false}}, `null`},
		{"e", Int32{v: sql.NullInt32{Int32: 0, Valid: true}}, `0`},
		{"f", Int32{v: sql.NullInt32{Int32: 1, Valid: true}}, `1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsequal(t, tt.in, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSuccessInt32Unmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
		want Int32
	}{
		{"a", `null`, Int32{}},
		{"b", `0`, NewInt32(0)},
		{"c", `1`, NewInt32(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Int32
			um(t, []byte(tt.in), &o)
			if !reflect.DeepEqual(tt.want, o) {
				t.Errorf(`want:%+v, got:%+v`, tt.want, o)
			}
		})
	}
}

func TestFailInt32Unmarshall(t *testing.T) {

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
			var o Int32
			if err := json.Unmarshal([]byte(tt.in), &o); err == nil {
				t.Error("err is nil")
			}
		})
	}
}
