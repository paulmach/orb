package geo

import "testing"

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

func BenchmarkPointUnmarshalWKB(b *testing.B) {
	p := NewPoint(0, 0)
	data := []uint8{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64}
	err := p.unmarshalWKB(data)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.unmarshalWKB(data)
	}
}

func BenchmarkPointSetScan(b *testing.B) {
	ps := NewPointSet()

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

func BenchmarkPointSetUnmarshalWKB(b *testing.B) {
	ps := NewPointSet()

	err := ps.unmarshalWKB(testPathWKB)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.unmarshalWKB(testPathWKB)
	}
}
