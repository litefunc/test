package main

import (
	"test/var/pkg1"
	"test/var/pkg2"
	"time"
)

func main() {

	t := time.NewTicker(time.Second)
	for range t.C {
		go func() {
			for i := 0; i < 10; i++ {
				pkg1.Log(i)
			}
		}()

		go func() {
			for i := 0; i < 10; i++ {
				pkg2.Log(i)
			}
		}()

	}
}
