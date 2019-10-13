package main

import (
	"fmt"
	"sync"
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
	rr()
	fmt.Println("finish")
}
