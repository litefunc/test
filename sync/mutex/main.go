package main

import (
	"cloud/lib/logger"
	"fmt"
	"sync"
	"time"
)

var m sync.RWMutex

func rr() {
	m.RLock()
	m.RLock()
	m.RUnlock()
	m.RUnlock()
}

// deadlock
func rw() {
	m.RLock()
	m.Lock()
	m.RUnlock()
	m.Unlock()
}

// deadlock
func wr() {
	m.Lock()
	m.RLock()
	m.Unlock()
	m.RUnlock()
}

func main() {
	// rr()
	mutex(1)
	mutex(2)
	mutex(3)
	fmt.Println("finish")
	var wc chan int
	<-wc
}

func mutex(n int) {
	var mu sync.Mutex
	for i := 0; i < 10; i++ {
		go func(i int) {
			mu.Lock()
			defer mu.Unlock()
			logger.Debug(n)
			time.Sleep(time.Second)

		}(i)
	}
}
