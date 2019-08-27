package main

import (
	"cloud/lib/logger"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		ch <- 1
		time.Sleep(time.Second * 2)
		close(ch)
	}()

	time.Sleep(time.Second)

	logger.Debug(<-ch)

	logger.Debug(<-ch)

	// for i := range ch {
	// 	logger.Debug(i)
	// }

}
