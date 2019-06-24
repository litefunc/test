package main

import "cloud/lib/logger"

type Allergen int

const (
	IgEggs         Allergen = 1 << iota // 1 << 0 which is 00000001
	IgChocolate                         // 1 << 1 which is 00000010
	IgNuts                              // 1 << 2 which is 00000100
	IgStrawberries                      // 1 << 3 which is 00001000
	IgShellfish                         // 1 << 4 which is 00010000
)

func main() {
	var a Allergen
	a = IgChocolate | IgNuts
	switch a {
	case IgChocolate | IgNuts:
		logger.Debug(a)
	case IgStrawberries | IgShellfish:
		logger.Debug(a, a)
	}

	logger.Debug(IgStrawberries & a)
	logger.Debug(IgNuts & a)
	logger.Debug(IgNuts&IgChocolate | IgNuts)
	logger.Debug(IgNuts | a)
}
