package main

import (
	"VodoPlay/logger"
	"encoding/json"
	"math"
	"time"
)

type A struct {
	A []int
	B []int `json:",nilasempty"`
	C []int
}

type B struct {
	A int `json:"a"`
	B *A  `json:"b"`
}

type C struct {
	F func() `json:"f"`
}

type Config struct {
	Auths           map[string]ConfigAuth `json:"auths"`
	json.RawMessage `json:",omitempty"`
}

type ConfigAuth struct {
	Auth            string `json:"auth"`
	json.RawMessage `json:",omitempty"`
}

func main() {
	var n interface{}
	var err error
	var by []byte

	by, err = json.Marshal(n)
	logger.Debugf(`%s`, by)
	logger.Error(err)

	by, err = json.Marshal(1)
	logger.Debugf(`%s`, by)
	logger.Error(err)

	var ch chan int
	by, err = json.Marshal(ch)
	logger.Debugf(`%s`, by)
	logger.Error(err)

	by, err = json.Marshal(math.Inf(1))
	logger.Debugf(`%s`, by)
	logger.Error(err)

	var p *struct{}
	by, err = json.Marshal(p)
	logger.Debugf(`%s`, by)
	logger.Error(err)

	by, err = json.Marshal(`{123`)
	if err != nil {
		logger.Error(err)
	}
	logger.Debugf(string(by))

	var i A
	i.C = make([]int, 0)
	by, _ = json.Marshal(i)
	logger.Debug(string(by))

	var i1 A
	json.Unmarshal(by, &i1)
	by, _ = json.Marshal(i1)
	logger.Debug(string(by))

	by, _ = json.Marshal(make([]int, 0))
	logger.Debug(string(by))

	is := make([]int, 0)
	json.Unmarshal([]byte("null"), &is)
	by, _ = json.Marshal(is)
	logger.Debug(string(by))

	type TestObject struct {
		Field1 string
		Field2 json.RawMessage
	}
	var data TestObject
	json.Unmarshal([]byte(`{"field1": "hello", "field2": [1,2,3]}`), &data)
	logger.Debug(data)
	by, _ = json.Marshal(data)
	logger.Debug(string(by))

	var b B
	logger.Debug(b)
	json.Unmarshal([]byte(`{"a": 1, "b": {}}`), &b)
	logger.Debug(b, *b.B)

	var c C
	by, err = json.Marshal(c)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(string(by))

	j(time.Hour * 365)

	m := map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}
	j(m)

	var cfg = Config{Auths: make(map[string]ConfigAuth)}
	cfg.Auths["a"] = ConfigAuth{Auth: "a"}
	j(cfg)

}

func j(o interface{}) {
	logger.Debug(o)
	by, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by))
}
