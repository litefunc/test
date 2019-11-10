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
}
