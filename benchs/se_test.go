package benchs

import "testing"

func BenchmarkSEAppend(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConcatAppend(100000)
	}
}

func BenchmarkSECopy(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConcatCopy(100000)
	}
}
func BenchmarkSEBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConcatBuilderPreGrow(100000)
	}
}
