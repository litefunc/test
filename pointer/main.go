package main

import "test/logger"

func main() {

	a := "a"
	b := &a
	c := &b
	logger.Debug(c, *c, b, &**c, **c)
}
