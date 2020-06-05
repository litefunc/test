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

	b := by("ab")
	logger.Debug(b)
	logger.Debug(tobyte(b))
	logger.Debug(tobyte(string(b)))

	p := &b
	logger.Debug(p)
	logger.Debug((*p)[0:0])
	logger.Debug(append((*p)[0:0], b...))
	*p = nil
	logger.Debug(b)

	b1 := make([]byte, 0, 0)
	b1 = by("ab")
	logger.Debug(b1)

	p = nil
	*p = nil
	logger.Debug(p)

}

func tobyte(o interface{}) bool {
	_, ok := o.([]byte)
	if !ok {
		return false
	}
	return true
}
