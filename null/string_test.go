package null

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringMarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   String
		want string
	}{
		{"a", String{}, `null`},
		{"b", NewString(""), `""`},
		{"c", NewString("s"), `"s"`},
		{"d", String{v: sql.NullString{String: "s", Valid: false}}, `null`},
		{"e", String{v: sql.NullString{String: "", Valid: true}}, `""`},
		{"f", String{v: sql.NullString{String: "s", Valid: true}}, `"s"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsequal(t, tt.in, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSuccessStringUnmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
		want String
	}{
		{"a", `null`, String{}},
		{"b", `""`, NewString("")},
		{"c", `"s"`, NewString("s")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o String
			um(t, []byte(tt.in), &o)
			if !reflect.DeepEqual(tt.want, o) {
				t.Errorf(`want:%+v, got:%+v`, tt.want, o)
			}
		})
	}
}

func TestFailStringUnmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
	}{
		{"unquoted", `s`},
		{"bool", `true`},
		{"int", `1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o String
			if err := json.Unmarshal([]byte(tt.in), &o); err == nil {
				t.Error("err is nil")
			}
		})
	}
}
