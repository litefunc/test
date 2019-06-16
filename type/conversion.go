package main

import "fmt"

type Ints []int

func (is Ints) Ints() []int {
	return []int{1, 2, 3}
}

func pins(is []int) {
	fmt.Println(is)
}

func pins1(is Ints) {
	fmt.Println(is.Ints())
}

type S struct{}

type Ss []S

func pss(is []S) {
	fmt.Println(is)
}

func pss1(is Ss) {
	fmt.Println(is)
}

func main() {

	var is []int
	pins(is)
	pins1(is)

	var is1 Ints
	pins(is1)
	pins1(is1)

	x := Ints(is)
	y := []int(is1)

	fmt.Println(x, y)

	var ss []S
	pss(ss)
	pss1(ss)

	var ss1 Ss
	pss(ss1)
	pss1(ss1)
}
