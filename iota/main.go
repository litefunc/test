package main

import "test/logger"

type Allergen int

// "<<" is interger operator, just like "+" "-" "*" "/", can be used anywhere, "iota" is only for constants, they are different concept

const (
	IgEggs         Allergen = 1 << iota // 00000001 ; iota=0 ; Allergen = 1 << iota = 1*1
	IgChocolate                         // 00000010 ; iota=1 ; Allergen = 1 << iota = 1*2
	IgNuts                              // 00000100 ; iota=2 ; Allergen = 1 << iota = 1*4
	IgStrawberries                      // 00001000 ; iota=3 ; Allergen = 1 << iota = 1*8
	IgShellfish                         // 00010000 ; iota=4 ; Allergen = 1 << iota = 1*16
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
