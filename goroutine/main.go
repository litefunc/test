package main

import (
	"fmt"
	"time"
)

func main() {
	var n int
	m := 10
	for i := 0; i < m; i++ {
		go func(i int) {
			fmt.Println(i)
			ticker := time.NewTicker(2 * time.Second)
			for t := range ticker.C {
				fmt.Println(i, t)
				n++
				if n == m {
					fmt.Println(n)
					n = 0
					fmt.Println()
				}
			}

			fmt.Println(i)
		}(i)
	}
	fmt.Println("end")
	var wc chan int
	<-wc
}
