package main

import "test/logger"

func main() {

	ch := make(chan int)
	// var ch chan int
	close(ch)

	for i := 0; i < 100; i++ {
		select {

		default:
			logger.Debug("default")

		case i := <-ch:
			logger.Debug(i)

		}
	}

	logger.Trace(<-ch)
	close(ch)
	ch <- 1
}
