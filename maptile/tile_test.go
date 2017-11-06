package maptile

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
)

func TestAt(t *testing.T) {
	tile := At(orb.NewPoint(0, 0), 28)
	if b := tile.Bound(); b.Top() != 0 || b.Left() != 0 {
		t.Errorf("incorrect tile bound: %v", b)
	}

	// specific case
	if tile := At(orb.NewPoint(-87.65005229999997, 41.850033), 20); tile.X != 268988 || tile.Y != 389836 {
		t.Errorf("projection incorrect: %v", tile)
	}

	if tile := At(orb.NewPoint(-87.65005229999997, 41.850033), 28); tile.X != 68861112 || tile.Y != 99798110 {
		t.Errorf("projection incorrect: %v", tile)
	}

	for _, city := range mercator.Cities {
		tile := At(orb.Point{city[1], city[0]}, 31)
		c := tile.Center()

		if math.Abs(c[1]-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", c[1], city[0])
		}

		if math.Abs(c[0]-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", c[0], city[1])
		}
	}

	// test polar regions
	if tile := At(orb.NewPoint(0, 89.9), 30); tile.Y != 0 {
		t.Errorf("top of the world error: %d != %d", tile.Y, 0)
	}

	if tile := At(orb.NewPoint(0, -89.9), 30); tile.Y != (1<<30)-1 {
		t.Errorf("bottom of the world error: %d != %d", tile.Y, (1<<30)-1)
	}
}

func TestTileQuadkey(t *testing.T) {
	// default level
	level := Zoom(30)
	for _, city := range mercator.Cities {
		tile := At(orb.Point{city[1], city[0]}, level)
		p := tile.Center()

		if math.Abs(p[1]-city[0]) > mercator.Epsilon {
			t.Errorf("latitude miss match: %f != %f", p[1], city[0])
		}

		if math.Abs(p[0]-city[1]) > mercator.Epsilon {
			t.Errorf("longitude miss match: %f != %f", p[0], city[1])
		}
	}
}

func TestTileBound(t *testing.T) {
	bound := Tile{7, 8, 9}.Bound()

	level := Zoom(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint32(5)

	// edges should be within the bound
	p := Tile{7<<factor + 1, 8<<factor + 1, level}.Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	p = Tile{7<<factor - 1, 8<<factor - 1, level}.Center()
	if bound.Contains(p) {
		t.Errorf("should not contain point")
	}

	p = Tile{8<<factor - 1, 9<<factor - 1, level}.Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	p = Tile{8<<factor + 1, 9<<factor + 1, level}.Center()
	if bound.Contains(p) {
		t.Errorf("should not contain point")
	}

	if b := New(0, 0, 0).Bound(); b != orb.NewBound(-180, 180, -85.05112877980659, 85.05112877980659) {
		t.Errorf("should be full earth, got %v", b)
	}

	if b := New(0, 0, 0).Bound(10); b != orb.NewBound(-180, 180, -85.05112877980659, 85.05112877980659) {
		t.Errorf("should be full earth, got %v", b)
	}
}

func TestFraction(t *testing.T) {
	p := Fraction(orb.NewPoint(-180, 0), 30)
	if p[0] != 0 {
		t.Errorf("should have left at zero: %f", p[0])
	}

	p = Fraction(orb.NewPoint(180, 0), 30)
	if p[0] != 0 {
		t.Errorf("should have right at zero: %f", p[0])
	}

	p = Fraction(orb.NewPoint(360, 0), 30)
	if p[0] != 1<<29 {
		t.Errorf("should have center: %f", p[0])
	}
}

func TestSharedParent(t *testing.T) {
	p := orb.NewPoint(-122.2711, 37.8044)
	one := At(p, 15)
	two := At(p, 15)

	expected := one

	one.Z = 25
	one.X = (one.X << 10) | 0x25A
	one.Y = (one.Y << 10) | 0x14B

	two.Z = 25
	two.X = (two.X << 10) | 0x15B
	two.Y = (two.Y << 10) | 0x26A

	if tile := one.SharedParent(two); tile != expected {
		t.Errorf("incorrect shared: %v != %v", tile, expected)
	}

	if tile := two.SharedParent(one); tile != expected {
		t.Errorf("incorrect shared: %v != %v", tile, expected)
	}
}

func BenchmarkSharedParent_SameZoom(b *testing.B) {
	p := orb.NewPoint(-122.2711, 37.8044)
	one := At(p, 10)
	two := At(p, 10)

	one.Z = 20
	one.X = (one.X << 10) | 0x25A
	one.Y = (one.X << 10) | 0x14B

	two.Z = 20
	two.X = (two.X << 10) | 0x15B
	two.Y = (two.X << 10) | 0x26A

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		one.SharedParent(two)
	}
}

func BenchmarkSharedParent_DifferentZoom(b *testing.B) {
	p := orb.NewPoint(-122.2711, 37.8044)
	one := At(p, 10)
	two := At(p, 10)

	one.Z = 20
	one.X = (one.X << 10) | 0x25A
	one.Y = (one.X << 10) | 0x14B

	two.Z = 18
	two.X = (two.X << 8) | 0x03B
	two.Y = (two.X << 8) | 0x0CA

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		one.SharedParent(two)
	}
}
