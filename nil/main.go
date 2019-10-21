package main

import (
	"fmt"
	"time"
)

var in interface{}

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
	time.Sleep(time.Second * 5)
}
