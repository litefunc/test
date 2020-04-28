package main

import (
	"fmt"

	"gotest.tools/assert/cmp"
)

func main() {
	fmt.Println(cmp.Equal(1, 2))
}
