package planar

import "testing"

func BenchmarkMultiPointCentroid(b *testing.B) {
	mp := append(NewMultiPoint(),
		NewPoint(0, 0),
		NewPoint(1, 1.5),
		NewPoint(2, 0),
		NewPoint(3, 1),
		NewPoint(3, 5),
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mp.Centroid()
	}
}
