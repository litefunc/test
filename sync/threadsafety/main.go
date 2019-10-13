package main

import (
	"fmt"
	"sync"
)

var n int

func main() {
	wg := new(sync.WaitGroup)

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			f(wg)
		}(i)
	}

	wg.Wait()
	fmt.Println(n)

}

func f(wg *sync.WaitGroup) {
	defer wg.Done()
	m := n + 1
	n = m
}
