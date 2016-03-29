package geo

import (
	"testing"

	"github.com/paulmach/orb/internal/wkb"
)

func BenchmarkPointScan(b *testing.B) {
	p := NewPoint(0, 0)
	data := []uint8{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64}
	err := p.Scan(data)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Scan(data)
	}
}

func BenchmarkPathScan(b *testing.B) {
	ps := NewPath()

	testPathWKB := wkb.PathTestCases[0].Data
	err := ps.Scan(testPathWKB)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Scan(testPathWKB)
	}
}
