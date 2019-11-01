package main

import (
	"cloud/lib/logger"
	"net/http"
	"sync"
)

func main() {

	var wc sync.WaitGroup
	for i := 0; i < 10; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()

			_, err := http.Get("http://localhost:8080/sleep")
			if err != nil {
				logger.Error(err)
			}

		}(i)
	}
	wc.Wait()

	for i := 0; i < 10; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()

			_, err := http.Get("http://localhost:8080/sleep1")
			if err != nil {
				logger.Error(err)
			}

		}(i)
	}
	wc.Wait()
}
