package geo_test

import (
	"testing"

	"github.com/paulmach/orb/geo"
)

func BenchmarkPointDistanceFrom(b *testing.B) {
	p1 := geo.NewPoint(-122.4167, 37.7833)
	p2 := geo.NewPoint(37.7833, -122.4167)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p1.DistanceFrom(p2)
	}
}

func BenchmarkPointDistanceFromHaversine(b *testing.B) {
	p1 := geo.NewPoint(-122.4167, 37.7833)
	p2 := geo.NewPoint(37.7833, -122.4167)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p1.DistanceFrom(p2, true)
	}
}

func BenchmarkPointQuadKey(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Quadkey(60)
	}
}

func BenchmarkPointQuadKeyString(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.QuadkeyString(60)
	}
}

func BenchmarkPointGeoHash(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GeoHash(12)
	}
}

func BenchmarkPointGeoHashInt64(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GeoHashInt64(60)
	}
}
