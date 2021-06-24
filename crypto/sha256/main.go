package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {

	sum := sha256.Sum256([]byte("hello world"))
	fmt.Printf("%x\n", sum)
	fmt.Println(sha256String("hello world"))
}

func sha256String(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}
