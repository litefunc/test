package main

import (
	"VodoPlay/logger"
	"fmt"
	"time"
)

var in interface{}

type A struct {
	A int
}

func get() int {
	time.Sleep(time.Second * 2)
	return 1
}

func main() {
	go func() {
		in = get()
	}()

	go func() {
		for {
			fmt.Println(in == nil, in)
			time.Sleep(time.Second)
		}
	}()

	var i []int

	i = nil
	logger.Debug(len(i))

	var a *A
	logger.Debug(*a)

	time.Sleep(time.Second * 5)
}
