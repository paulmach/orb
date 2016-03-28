package planar

import "testing"

func BenchmarkPointSetCentroid(b *testing.B) {
	ps := PointSet{}
	ps = append(ps,
		Point{0, 0},
		Point{1, 1.5},
		Point{2, 0},
		Point{3, 1},
		Point{3, 5},
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Centroid()
	}
}
