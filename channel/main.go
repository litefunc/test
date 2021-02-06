package main

import (
	"test/logger"
)

func main() {
	ch := make(chan int)

	go func() {
		ch <- 1
		close(ch)

	}()

	logger.Debug(<-ch)
	logger.Debug(<-ch)
	logger.Debug(<-ch)

	for i := range ch {
		logger.Debug(i)
	}
	ch <- 2

	// ch1 := make(chan []*int)
	// var is []*int
	// for i := 0; i < 10; i++ {
	// 	is = append(is, &i)
	// }

	// go func() {
	// 	ch1 <- is
	// 	is = nil
	// }()
	// time.Sleep(time.Second)
	// got := <-ch1
	// logger.Debug(len(got), got)
}
