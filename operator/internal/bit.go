package internal

import (
	"math"
	"test/logger"
)

func Bit() {
	for i := 0; i <= 10; i++ {
		logger.Tracef(`%d:%b`, i, i)
		logger.Debugf(`%d:%b`, i, -i)
		logger.Infof(`%d:%b`, i, uint64(i))
	}
}

func Uint32() {
	for i := 0; i <= 33; i++ {
		f := math.Pow(2, float64(i))
		u := uint32(f)
		logger.Tracef(`%d,%d:%b`, i, u, u)
	}
}

func Left() {
	n := uint32(3)
	n1 := uint32(2147483647)
	n2 := uint32(2147483648)

	for i := 0; i <= 10; i++ {
		logger.Tracef(`%d:%b`, n<<i, n<<i)
		logger.Debugf(`%d:%b`, n1<<i, n1<<i)
		logger.Infof(`%d:%b`, n2<<i, n2<<i)
	}
}
