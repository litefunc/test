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

type ints []int

func (list ints) add(n int) {
	for i, v := range list {
		list[i] = v + n
	}
}

func add(list ints, n int) {
	for i, v := range list {
		list[i] = v + n
	}
}

func (list ints) append(n int) {
	list = append(list, n)
}

type S2 struct {
	A []int
	B []int
	C int
}

func (s S2) setA(a []int) {
	logger.Debugf(`%p %p %p %p`, &s, &s.A, &s.B, &s.C)
	s.A = a
}

func main() {

	// s1 := make([]int, 2, 3)
	// logger.Debug(s1, len(s1), cap(s1))
	// s1[0] = 0
	// s1[1] = 1
	// s1 = append(s1, 2, 3)
	// logger.Debug(s1, len(s1), cap(s1))
	// s1 = append(s1, 4, 5, 6)
	// logger.Debug(s1, len(s1), cap(s1))

	// var s2 []int
	// logger.Debug(cap(s2))

	// var s3 []int
	// s2 = nil
	// s3 = nil
	// s4 := append(s2, s3...)
	// logger.Debug(s4)

	// is := ints{0, 1, 2}
	// is.add(1)
	// logger.Debug(is)
	// add(is, 1)
	// logger.Debug(is)
	// is.append(5)
	// logger.Debug(is)

	// s5 := S2{A: []int{1, 2, 3}}
	// logger.Debugf(`%p %p %p %p`, &s5, &s5.A, &s5.B, &s5.C)
	// s5.setA([]int{2, 3, 4})
	// logger.Debug(s5)

	var ss []*S1
	var s *S1
	s = &S1{A: 1}
	ss = append(ss, s)
	s = &S1{A: 2}
	ss = append(ss, s)
	logger.Debug(ss)
	logger.Debug(ss[0], ss[1])
}
