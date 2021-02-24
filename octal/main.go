package main

import "test/logger"

func main() {
	var i int
	// integer begins with "0" is octal, otherwise is decimal
	i = 024 //Equivalent to 20 in decimal
	logger.Trace(i + 10)
	logger.Trace(777)
	logger.Trace(0777, 7+7*8+7*64)
	logger.Trace(07777, 7+7*8+7*64+7*512)
	logger.Trace(0011&0101, 0001, 1)
	logger.Trace(0011|0101, 0111, 1+1*8+1*64)
	logger.Trace(01011|0101, 01111, 1+1*8+1*64+1*512)

}
