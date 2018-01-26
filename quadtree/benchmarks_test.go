package quadtree

import (
	"math"
	"math/rand"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

func BenchmarkAdd(b *testing.B) {
	r := rand.New(rand.NewSource(22))
	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qt.Add(orb.Point{r.Float64(), r.Float64()})
	}
}

func BenchmarkRandomFind1000(b *testing.B) {
	r := rand.New(rand.NewSource(42))
	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})

	for i := 0; i < 1000; i++ {
		qt.Add(orb.Point{r.Float64(), r.Float64()})
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qt.Find(orb.Point{r.Float64(), r.Float64()})
	}
}

func BenchmarkRandomFind1000Naive(b *testing.B) {
	r := rand.New(rand.NewSource(42))

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	points := []orb.Point{}

	for i := 0; i < 1000; i++ {
		p := orb.Point{r.Float64(), r.Float64()}

		qt.Add(p)
		points = append(points, p)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		looking := orb.Point{r.Float64(), r.Float64()}

		min := math.MaxFloat64
		var best orb.Point
		for _, p := range points {
			if d := planar.DistanceSquared(looking, p); d < min {
				min = d
				best = p
			}
		}

		_ = best
	}
}

func BenchmarkRandomInBound1000(b *testing.B) {
	r := rand.New(rand.NewSource(43))

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		qt.Add(orb.Point{r.Float64(), r.Float64()})
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := orb.Point{r.Float64(), r.Float64()}
		qt.InBound(nil, p.Bound().Pad(0.1))
	}
}

func BenchmarkRandomInBound1000Naive(b *testing.B) {
	r := rand.New(rand.NewSource(43))

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	points := []orb.Point{}

	for i := 0; i < 1000; i++ {
		p := orb.Point{r.Float64(), r.Float64()}

		qt.Add(p)
		points = append(points, p)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var near []orb.Point
	for i := 0; i < b.N; i++ {
		p := orb.Point{r.Float64(), r.Float64()}
		b := orb.Bound{Min: p, Max: p}
		b = b.Pad(0.1)

		near = near[:0]
		for _, p := range points {
			if b.Contains(p) {
				near = append(near, p)
			}
		}

		_ = len(near)
	}
}

func BenchmarkRandomInBound1000Buf(b *testing.B) {
	r := rand.New(rand.NewSource(43))

	qt := New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		qt.Add(orb.Point{r.Float64(), r.Float64()})
	}

	var buf []orb.Pointer
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := orb.Point{r.Float64(), r.Float64()}
		buf = qt.InBound(buf, p.Bound().Pad(0.1))
	}
}
