package internal

import "LocalServer/lib/logger"

type sa struct {
	a int
	b []int
}

func Copy() {

	a := sa{a: 1, b: []int{1, 2, 3}}
	b := a
	b.a = 2
	b.b = []int{4, 5, 6}
	logger.Debug(a)
	logger.Debug(b)

}
