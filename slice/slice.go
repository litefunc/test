package main

import "cloud/lib/logger"

type S1 struct {
	A int
}

func Set1(s []int, i int) {
	s[0] = i
	logger.Debugf(`%p`, s)
}

func SetA(s []S1, i, v int) {
	s[i].A = v
}

func main() {

	var s1, s2 []int
	var ch chan int
	var m map[int]int
	var i, j interface{}

	logger.Debug(s1, s2, ch, m, i, j)
	logger.Debugf(`%p, %p, %p, %p, %p, %p`, s1, s2, ch, m, &i, &j)

	s1 = []int{1, 2, 3}
	s2 = s1
	Set1(s1, 2)
	logger.Debug(s1)

	ch = make(chan int)
	m = make(map[int]int)
	i = 1
	j = "s"

	logger.Debug(s1, s2, ch, m, i, j)
	logger.Debugf(`%p, %p, %p, %p, %p, %p`, s1, s2, ch, m, &i, &j)

	var s11, s12 S1

	s1s := []S1{s11, s12}
	SetA(s1s, 0, 1)
	logger.Debug(s1s)
}
