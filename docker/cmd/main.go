package main

import "test/docker/cmd/internal"

func main() {
	internal.Ping()
	internal.ImageList()
	// internal.ContainerList()

}
