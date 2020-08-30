package main

import (
	"log"
	"os"
	"syscall"
	"test/log/internal"
)

func main() {

	// internal.SetStderr("logfile.txt")

	// panic(3)

	// defer func() {
	// 	if x := recover(); x != nil {
	// 		// recovering from a panic; x contains whatever was passed to panic()
	// 		log.Printf("run time panic: %v", x)

	// 		// if you just want to log the panic, panic again
	// 		panic(x)
	// 	}
	// }()

	for 
}

func p1() {
	var ch chan int
	go p2()
	<-ch
}

func p2() {
	panic(2)
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}
