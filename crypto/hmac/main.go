package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {

	secret := "mysecret"
	data := "data"
	fmt.Printf("Secret: %s Data: %s\n", secret, data)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	fmt.Println("Result: " + sha)

	fmt.Println(hmacString(data, secret))
}

func hmacString(data, secret string) string {

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}
