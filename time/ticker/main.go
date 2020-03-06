package main

import (
	"cloud/lib/logger"
	"time"
)

func main() {

	n := 1

	ticker := time.NewTicker(time.Second * time.Duration(n))

	for {
		select {
		case c := <-ticker.C:
			logger.Debug(n, c.UTC())
			if n == 1 {
				n++
				ticker.Stop()
				ticker = time.NewTicker(time.Second * time.Duration(n))

			} else {
				n--
				ticker.Stop()
				ticker = time.NewTicker(time.Second * time.Duration(n))
			}

		}
	}
	logger.Debug("done")
}

// func main() {

// 	n := 1

// 	ticker := time.NewTicker(time.Second * time.Duration(n))

// 	for c := range ticker.C {

// 		logger.Debug(n, c.UTC())
// 		ticker.Stop()

// 	}
// 	logger.Debug("done")
// }
