package main

import (
	"os"
	"path"
)

func main() {

	p := path.Join(os.Getenv("GOPATH"), "src/test/os/append/test.txt")

	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, s := range []string{"a", "b", "c", "d"} {
		if _, err = f.WriteString(s); err != nil {
			panic(err)
		}
	}

}
