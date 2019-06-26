package main

import (
	"cloud/lib/logger"
	"context"
	"time"
)

var n int

func g(i int, ctx context.Context) {

	for {
		select {
		case s := <-ctx.Done():
			logger.Warn(s)
			return
		default:
			logger.Debug(i, n)
			n++
			time.Sleep(time.Second)
		}
	}

}

func main() {

	var wc chan int
	// wc := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for i := 0; i <= 3; i++ {
		go g(i, ctx)
		time.Sleep(time.Second)
	}

	logger.Debug(<-wc)
}
