package benchs

import "testing"

func BenchmarkCompareBytes(b *testing.B) {
	a := []byte("hello")
	for i := 0; i < b.N; i++ {
		CompareV1(a)
	}
}

func BenchmarkCompareStrings(b *testing.B) {
	a := []byte("hello")
	for i := 0; i < b.N; i++ {
		CompareV2(a)
	}
}
