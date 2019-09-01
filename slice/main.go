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

	s1 := make([]int, 2, 3)
	logger.Debug(s1, len(s1), cap(s1))
	s1[0] = 0
	s1[1] = 1
	s1 = append(s1, 2, 3)
	logger.Debug(s1, len(s1), cap(s1))
	s1 = append(s1, 4, 5, 6)
	logger.Debug(s1, len(s1), cap(s1))

	var s2 []int
	logger.Debug(cap(s2))
}
