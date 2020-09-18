package main

import (
	"LocalServer/logger"
	"fmt"
	"log"
	"os"
)

func main() {

	// total := 179496183
	total := 0

	// p := os.Getenv("GOPATH") + "/src/test/os/size/FLUID_165099_49099M_022044.mp4"
	p := os.Getenv("GOPATH") + "/src/test/os/size/empty.txt"

	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	size := fi.Size()

	if size == 0 {
		size = 1
	}
	fmt.Println(size)

	var percent float64 = float64(size) / float64(total) * 100
	logger.Debug(percent, percent > 0)
	fmt.Printf("%.0f", percent)
	fmt.Println("%")

}
