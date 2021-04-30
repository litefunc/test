package main

import (
	"sync"
	"test/logger"
	"time"
)

func sleep(i int64, wc *sync.WaitGroup) {
	defer wc.Done()
	t := time.Now()
	time.Sleep(time.Second * 5)
	logger.Trace(i, time.Now().Sub(t))

}

func tick(i int64, wc *sync.WaitGroup) {
	defer wc.Done()
	ticker := time.NewTicker(time.Second * 5)
	for i := range ticker.C {
		logger.Trace(i)
	}

}

func main() {

	// fmt.Println(runtime.NumCPU())
	// fmt.Println(runtime.GOMAXPROCS(0))
	// runtime.GOMAXPROCS(1)
	// fmt.Println(runtime.NumCPU())

	// fmt.Println("start", runtime.GOMAXPROCS(1), runtime.NumGoroutine())

	// var n int
	// m := 10
	// for i := 0; i < m; i++ {
	// 	go func(i int) {
	// 		fmt.Println(i)
	// 		ticker := time.NewTicker(2 * time.Second)
	// 		for t := range ticker.C {
	// 			fmt.Println(i, t)
	// 			n++
	// 			if n == m {
	// 				fmt.Println(n)
	// 				n = 0
	// 				fmt.Println(runtime.NumGoroutine())
	// 				fmt.Println()
	// 			}
	// 		}

	// 		fmt.Println(i)
	// 	}(i)
	// }
	// fmt.Println("end")
	// var wc chan int
	// <-wc

	t := time.Now()
	wc := &sync.WaitGroup{}
	wc.Add(2)
	go tick(1, wc)
	go tick(2, wc)
	wc.Wait()
	logger.Debug(time.Now().Sub(t))
}
