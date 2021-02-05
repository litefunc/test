package null

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
)

func TestFloat64Marshall(t *testing.T) {

	var tests = []struct {
		name string
		in   Float64
		want string
	}{
		{"a", Float64{}, `null`},
		{"b", NewFloat64(0), `0`},
		{"c", NewFloat64(1), `1`},
		{"d", Float64{v: sql.NullFloat64{Float64: 1, Valid: false}}, `null`},
		{"e", Float64{v: sql.NullFloat64{Float64: 0, Valid: true}}, `0`},
		{"f", Float64{v: sql.NullFloat64{Float64: 1, Valid: true}}, `1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsequal(t, tt.in, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSuccessFloat64Unmarshall(t *testing.T) {

	var tests = []struct {
		name string
		in   string
		want Float64
	}{
		{"a", `null`, Float64{}},
		{"b", `0`, NewFloat64(0)},
		{"c", `1`, NewFloat64(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Float64
			um(t, []byte(tt.in), &o)
			if !reflect.DeepEqual(tt.want, o) {
				t.Errorf(`want:%+v, got:%+v`, tt.want, o)
			}
		})
	}
}

func TestFailFloat64Unmarshall(t *testing.T) {

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
			var o Float64
			if err := json.Unmarshal([]byte(tt.in), &o); err == nil {
				t.Error("err is nil")
			}
		})
	}
}
