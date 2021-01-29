package main

import (
	"fmt"
	"test/logger"
)

type Ints []int

func (is Ints) Ints() []int {
	return []int{1, 2, 3}
}

func pins(is []int) {
	logger.Debug(is)
}

func pins1(is Ints) {
	logger.Debug(is.Ints())
}

type S struct{}

type Ss []S

func pss(is []S) {
	fmt.Println(is)
}

func pss1(is Ss) {
	fmt.Println(is)
}

type I int

func pi(i int) {
	fmt.Println(i)
}

func pI(i I) {
	fmt.Println(i)
}

func main() {

	var is = []int{1, 2, 3}
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

	var i0 I
	var i1 int

	pI(i0)
	// pI(i1)
	pI(I(i1))
	pI(1)

	// pi(i0)
	pi(i1)
	pi(int(i0))
	pi(1)
}
