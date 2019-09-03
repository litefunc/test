package slice

import "testing"

func BenchmarkSlice01(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var list []int
		for i := 0; i < 1000; i++ {
			list = append(list, i)
		}
	}
}

func BenchmarkSlice02(b *testing.B) {

	for i := 0; i < b.N; i++ {
		list := make([]int, 0, 1000)
		for i := 0; i < 1000; i++ {
			list = append(list, i)
		}
	}
}

func BenchmarkSlice03(b *testing.B) {

	for i := 0; i < b.N; i++ {
		list := make([]int, 1000, 1000)
		for i := 0; i < 1000; i++ {
			list = append(list, i)
		}
	}
}

func BenchmarkSlice05(b *testing.B) {

	for i := 0; i < b.N; i++ {
		list := make([]int, 1000, 1000)
		for i := 0; i < 1000; i++ {
			list[0] = i
		}
	}
}

func BenchmarkSlice06(b *testing.B) {

	for i := 0; i < b.N; i++ {
		list := make([]int, 1000)
		for i := 0; i < 1000; i++ {
			list[0] = i
		}
	}
}

func BenchmarkArray01(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var list [1000]int
		for i := 0; i < 1000; i++ {
			list[0] = i
		}
	}
}

func BenchmarkArray02(b *testing.B) {

	for i := 0; i < b.N; i++ {
		list := [1000]int{}
		for i := 0; i < 1000; i++ {
			list[0] = i
		}
	}
}
