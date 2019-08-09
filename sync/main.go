package main

import (
	"cloud/lib/logger"
	"sync"
	"time"
)

func lock(mu *sync.Mutex, wg *sync.WaitGroup, i int) {
	mu.Lock()
	logger.Debug("lock", i)
	defer mu.Unlock()
	defer wg.Done()
	time.Sleep(time.Second)
	logger.Debug("unlock", i)
}

func main() {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	mutex.Unlock()

	var wg = &sync.WaitGroup{}
	for r := 0; r < 100; r++ {
		wg.Add(1)
		go lock(mutex, wg, r)
	}
	wg.Wait()
}
