package main

import (
	"cloud/lib/logger"
	"time"
)

func main() {
	// ch := make(chan int)

	// go func() {
	// 	ch <- 1
	// 	time.Sleep(time.Second * 2)
	// 	close(ch)
	// }()

	// time.Sleep(time.Second)

	// logger.Debug(<-ch)

	// logger.Debug(<-ch)

	// for i := range ch {
	// 	logger.Debug(i)
	// }
	ch1 := make(chan []*int)
	var is []*int
	for i := 0; i < 10; i++ {
		is = append(is, &i)
	}

	go func() {
		ch1 <- is
		is = nil
	}()
	time.Sleep(time.Second)
	got := <-ch1
	logger.Debug(len(got), got)
}
