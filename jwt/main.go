package main

import (
	"VodoPlay/logger"
	"test/jwt/internal"
	"time"
)

func main() {

	s1 := "!qaz2wsx"
	s2 := "test2"
	tk1 := internal.Gen("1", s1, time.Hour*24*365)
	tk2 := internal.Gen("2", s2, time.Second)
	tk3 := internal.Gen("3", s1, time.Second)
	tk4 := internal.Gen("4", s2, time.Second)

	for i, v := range []string{tk1, tk2, tk3, tk4} {
		logger.Debug(i, v)
	}

	time.Sleep(time.Second * 2)

	internal.Validate(1, s1, tk1)
	internal.Validate(2, s1, tk2)
	internal.Validate(3, s2, tk1)
	internal.Validate(4, s2, tk2)

	internal.Validate(5, s1, tk3)
	internal.Validate(6, s1, tk4)
	internal.Validate(7, s2, tk3)
	internal.Validate(8, s2, tk4)
}
