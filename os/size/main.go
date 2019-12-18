package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	total := 179496183

	p := os.Getenv("GOPATH") + "/src/test/os/size/FLUID_165099_49099M_022044.mp4"

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

	fmt.Printf("%.0f", percent)
	fmt.Println("%")

}
