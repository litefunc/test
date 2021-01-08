package main

import (
	"VodoPlay/logger"
	"encoding/json"
)

func main() {

	var m map[int]int

	um([]byte(`{}`), &m)
	m[1] = 1

	m1 := m

	m1[1] = 2
	logger.Debug(m, m1)

	ms := M{m: m1}
	m3 := ms.M()
	m3[1] = 3
	logger.Debug(m, m1, m3)
}

type M struct {
	m map[int]int
}

func (rec M) M() map[int]int {
	return rec.m
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

func um(data []byte, o interface{}) {
	logger.Debug(string(data))
	logger.Debug(o)

	if err := json.Unmarshal(data, o); err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(o)
}
