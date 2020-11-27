package main

import (
	"VodoPlay/logger"
	"encoding/json"
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

func main() {

	um([]byte{}, &A{})
	um([]byte{}, []int{})
	um([]byte(`[1,2]`), []int{})
	um([]byte(`[1,2]`), &[]int{})
	um([]byte(`["1","2"]`), &[]int{})
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
