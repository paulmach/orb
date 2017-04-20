package planar_test

import (
	"testing"

	"github.com/paulmach/orb/planar"
)

func BenchmarkPointDistanceFrom(b *testing.B) {
	p1 := planar.NewPoint(-122.4167, 37.7833)
	p2 := planar.NewPoint(37.7833, -122.4167)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p1.DistanceFrom(p2)
	}
}

func BenchmarkPointDistanceFromSquared(b *testing.B) {
	p1 := planar.NewPoint(-122.4167, 37.7833)
	p2 := planar.NewPoint(37.7833, -122.4167)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p1.DistanceFromSquared(p2)
	}
}
