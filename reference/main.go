package main

import "cloud/lib/logger"

func sliceAppend(s []int, i int) {
	s = append(s, i)
}

func main() {

	var is []int

	sliceAppend(is, 1)

	logger.Debug(is)
}
