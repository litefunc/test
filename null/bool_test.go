package null

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
)

func TestBoolMarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   Bool
		want string
	}{
		{"a", Bool{}, `null`},
		{"b", NewBool(false), `false`},
		{"c", NewBool(true), `true`},
		{"d", Bool{v: sql.NullBool{Bool: true, Valid: false}}, `null`},
		{"e", Bool{v: sql.NullBool{Bool: true, Valid: true}}, `true`},
		{"f", Bool{v: sql.NullBool{Bool: false, Valid: true}}, `false`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsequal(t, tt.in, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSuccessBoolUnmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
		want Bool
	}{
		{"a", `null`, Bool{}},
		{"b", `false`, NewBool(false)},
		{"c", `true`, NewBool(true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Bool
			um(t, []byte(tt.in), &o)
			if !reflect.DeepEqual(tt.want, o) {
				t.Errorf(`want:%+v, got:%+v`, tt.want, o)
			}
		})
	}
}

func TestFailBoolUnmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
	}{
		{"string", `"s"`},
		{"string-1", `"true"`},
		{"int", `1`},
		{"time", `"0001-01-01T00:00:00Z"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Bool
			if err := json.Unmarshal([]byte(tt.in), &o); err == nil {
				t.Error("err is nil")
			}
		})
	}
}
