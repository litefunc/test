package main

import (
	"cloud/lib/logger"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	// cli := hello.NewClient("localhost", 8889)

	// for i := 0; i < 30; i++ {
	// 	re, err := cli.Hello(&pb.Empty{})
	// 	if err != nil {
	// 		return
	// 	}
	// 	logger.Debug(i+1, re.Id)
	// }

	start := time.Now()
	var wc sync.WaitGroup
	for i := 0; i < 30; i++ {
		wc.Add(1)
		go func(i int) {
			defer wc.Done()
			res, err := http.Get("http://localhost:8080/wait")
			if err != nil {
				logger.Error(err)
				return
			}
			by, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logger.Error(err)
				return
			}
			logger.Debug(i+1, string(by))
		}(i)

	}
	wc.Wait()
	logger.Debug("end:", time.Now().Sub(start).Seconds())
}
