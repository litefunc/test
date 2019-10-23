package main

import (
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func main() {

	f, err := os.Open("MTV.png")
	if err != nil {
		log.Fatal(err)
	}

	// This example uses png.Decode which can only decode PNG images.
	// Consider using the general image.Decode as it can sniff and decode any registered image format.
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	bound := img.Bounds()
	log.Println(bound.Dx(), bound.Dy())

	m := resize.Resize(1744, 454, img, resize.Lanczos3)

	out, err := os.Create("test_resized.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	if err := png.Encode(out, m); err != nil {
		log.Fatal(err)
	}
}
