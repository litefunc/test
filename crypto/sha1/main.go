package main

import (
	"LocalServer/logger"
	"crypto/sha1"
	"encoding/hex"
)

func main() {

	s := "sha1 this string"
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))

	logger.Debug(s, sha1_hash)
}
