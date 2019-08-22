package main

import "fmt"

func p() int {
	defer func() int {
		fmt.Println("defer caller")
		if err := recover(); err != nil {
			fmt.Println("recover success.")
			return -1
		}
		return 2
	}()

	panic("p")
	return 1
}

func main() {

	fmt.Println(p())
}
