package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"test/logger"
)

func main() {
	bytes := make([]byte, 16) //generate a random 32 byte key for AES-128
	if _, err := rand.Read(bytes); err != nil {
		logger.Fatal(err)
	}
	logger.Debug(string(bytes))
	logger.Debug(hex.EncodeToString(bytes), fmt.Sprintf("%x", bytes))
}
