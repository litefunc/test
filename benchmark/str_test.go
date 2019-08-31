package benchmark

import (
	"fmt"
	"strconv"
	"testing"
)

func printInt2String01(num int) string {
	return fmt.Sprintf("%d", num)
}

func printInt2String02(num int64) string {
	return strconv.FormatInt(num, 10)
}
func printInt2String03(num int) string {
	return strconv.Itoa(num)
}

func BenchmarkPrintInt2String01(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printInt2String01(100)
	}
}

func BenchmarkPrintInt2String02(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printInt2String02(int64(100))
	}
}

func BenchmarkPrintInt2String03(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printInt2String03(100)
	}
}
