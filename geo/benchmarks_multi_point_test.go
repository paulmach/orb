package geo

import "testing"

func BenchmarkPointSetCentroid(b *testing.B) {
	mp := append(NewMultiPoint(),
		NewPoint(-188.1298828125, -33.97980872872456),
		NewPoint(-186.1083984375, -38.54816542304658),
		NewPoint(-194.8974609375, -46.10370875598026),
		NewPoint(-192.1728515625, -47.8721439688873),
		NewPoint(-179.7802734375, -37.30027528134431),
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mp.Centroid()
	}
}
