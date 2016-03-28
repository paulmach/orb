package planar_test

import (
	"testing"

	planar "."
)

func BenchmarkVectorNormalize(b *testing.B) {
	v := planar.NewVector(5, 6)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Normalize()
	}
}
