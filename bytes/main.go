package main

import (
	"cloud/lib/logger"
)

func by(s string) []byte {
	return []byte(s)
}

func st(is ...uint8) string {
	var b []byte
	// for _, v := range is {
	// 	b = append(b, v)
	// }
	return string(append(b, is...))
}

func ascii() {
	for i := 0; i < 128; i++ {
		logger.Debug(i, st(uint8(i)))
	}
}

func main() {

	ascii()

	logger.Debug(by("\n"))
	logger.Debug(by("\t"))

}
