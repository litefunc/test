package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	Filesize uint64 `json:"filesize"`
	Filename string `json:"filename"`
}

type B struct {
	filesize uint64 `json:"filesize"`
	filename string `json:"filename"`
}

type c struct {
	Filesize uint64 `json:"filesize"`
	Filename string `json:"filename"`
}
type d struct {
	filesize uint64 `json:"filesize"`
	filename string `json:"filename"`
}

var (
	byt = []byte(`{"filesize":123,"filename":"abc"}`)

	sa = A{
		Filesize: 123,
		Filename: "abc",
	}
	sb = B{
		filesize: 123,
		filename: "abc",
	}
	sc = c{
		Filesize: 123,
		Filename: "abc",
	}
	sd = d{
		filesize: 123,
		filename: "abc",
	}
)

func p(o interface{}) {

	by, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(by))
}

func u(o interface{}) {

	if err := json.Unmarshal(byt, o); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(o)
}

func main() {

	p(sa)
	p(sb)
	p(sc)
	p(sd)

	u(&A{})
	u(&B{})
	u(&c{})
	u(&d{})
}
