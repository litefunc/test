package main

import (
	"VodoPlay/logger"
	"flag"
)

func main() {

	new := flag.Bool("new", false, "add this device to database")
	name := flag.String("name", "", "device name")

	flag.Parse()

	logger.Debug(isFlagPassed("new"))
	logger.Debug(isFlagPassed("name"))

	if *new == true {
		return
	}

	if *name == "" {
		return
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
