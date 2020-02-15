package main

import "cloud/lib/logger"

type S struct {
	A int
}

func main() {

	s := new(S)
	logger.Debug(s)

	var s1 *S
	logger.Debug(s1)

	s2 := new(S)

	logger.Debug(s2)
	logger.Debug(s == s2)

	var i int
	var p *int
	logger.Debug(p, &p, &i)
	p = &i
	logger.Debug(p, *p, &p, &i)
	i = 1
	logger.Debug(p, *p, &p, &i)

	p1 := &i
	i = 2
	logger.Debug(p1, &p1, &i)
	sl := []int{0, 1, 2}
	logger.Debug(&sl, &sl[0], &sl[1])
	sl[0] = 1
	logger.Debug(&sl, &sl[0], &sl[1])

	sl1 := &sl
	logger.Debug(&sl1)

}
