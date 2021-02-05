package internal

import (
	"LocalServer/lib/logger"
	"encoding/json"
	"time"
)

func Null() {
	um([]byte(`null`))
	uma([]byte(`null`))
	uma([]byte(`{"a": 1, "b": 2}`))
	umb([]byte(`null`))
	umb([]byte(`{"a":  null, "b":"foo"}`))
	umb([]byte(`{"a": {"a": 1, "b": 2}, "b":"foo"}`))
}

func NullString() {
	um([]byte(`"null"`))
}

func Empty() {
	um([]byte(``))
}

func EmptyString() {
	um([]byte(`""`))
}

func um(by []byte) {
	logger.Info(by, string(by))

	var s string
	if err := json.Unmarshal(by, &s); err != nil {
		logger.Error(err)
	}
	logger.Debug("string:", s)

	var b bool
	if err := json.Unmarshal(by, &b); err != nil {
		logger.Error(err)
	}
	logger.Debug("bool:", b)

	var i int
	if err := json.Unmarshal(by, &i); err != nil {
		logger.Error(err)
	}
	logger.Debug("int:", i)

	var u uint64
	if err := json.Unmarshal(by, &u); err != nil {
		logger.Error(err)
	}
	logger.Debug("uint64:", u)

	var t time.Time
	if err := json.Unmarshal(by, &t); err != nil {
		logger.Error(err)
	}
	logger.Debug("time:", t)
}

func uma(by []byte) {
	logger.Info(string(by))

	var a A
	if err := json.Unmarshal(by, &a); err != nil {
		logger.Error(err)
	}
	b, err := json.Marshal(a)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("a:", string(b))
}

type A struct {
	valid bool
	a
}

type a struct {
	A int
	B int
}

func (rec A) MarshalJSON() ([]byte, error) {
	if !rec.valid {
		return []byte(`null`), nil
	}
	return json.Marshal(rec.a)
}

func (rec *A) UnmarshalJSON(b []byte) error {

	if string(b) == `null` {
		rec.valid = false
		return nil
	}

	err := json.Unmarshal(b, &rec.a)

	rec.valid = (err == nil)
	return err
}

type B struct {
	A A      `json:"a"`
	B string `json:"b"`
}

func umb(by []byte) {
	logger.Info(string(by))

	var b B
	if err := json.Unmarshal(by, &b); err != nil {
		logger.Error(err)
	}
	b1, err := json.Marshal(b)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("b:", string(b1))
}
