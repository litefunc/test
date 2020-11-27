package main

import (
	"LocalServer/lib/logger"
	"test/color/internal"
)

func main() {

	// for _, v := range []int{30, 31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96, 97} {
	// 	internal.Log(v, "[Debug] abcdrfg122323123123")
	// }

	// for _, v := range []int{40, 41, 42, 43, 44, 45, 46, 47, 100, 101, 102, 103, 104, 105, 106, 107} {
	// 	internal.Bg(v, "[Debug] abcdrfg122323123123")
	// }

	// for _, v := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
	// 	internal.B(v, "[Debug] abcdrfg122323123123")
	// }

	c := internal.NewColor(internal.FgBrBlue, internal.BgBrGreen, internal.CrossedOut)
	for i := 0; i < 1; i++ {
		c.Println("asdasdasdsad")
	}

	logger.HTTP(123)
	logger.Trace(123)
	logger.Debug(123)
	logger.Info(123)
	logger.Warn(123)
	logger.Error(123)
	// logger.Panic(123)
	logger.Fatal(123)
}
