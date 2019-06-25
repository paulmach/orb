package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	orb "github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// See https://github.com/mbostock/topojson/blob/master/test/topology/join-test.js

// join the returned hashmap has true for junction points
func TestHasJunctions(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", orb.LineString{
			orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("ab", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{0, 0}))
}

// join the returned hashmap has undefined for non-junction points
func TestNonJunctions(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", orb.LineString{
			orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("ab", orb.LineString{
			orb.Point{0, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.False(junctions.Has([]float64{1, 0}))
}

// join exact duplicate lines ABC & ABC have junctions at their end points
func TestJoinDuplicate(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("abc2", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 2)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join reversed duplicate lines ABC & CBA have junctions at their end points"
func TestJoinReversedDuplicate(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("cba", orb.LineString{
			orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 2)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join exact duplicate rings ABCA & ABCA have no junctions
func TestJoinDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("abca2", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join reversed duplicate rings ACBA & ABCA have no junctions
func TestJoinReversedDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("acba", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join rotated duplicate rings BCAB & ABCA have no junctions
func TestJoinRotatedDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("bcab", orb.Polygon{
			orb.Ring{
				orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0}, orb.Point{1, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join ring ABCA & line ABCA have a junction at A
func TestJoinRingLine(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcaLine", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("abcaPolygon", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{0, 0}))
}

// join ring BCAB & line ABCA have a junction at A
func TestJoinLineRingReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcaLine", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("bcabPolygon", orb.Polygon{
			orb.Ring{
				orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0}, orb.Point{1, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{0, 0}))
}

// join ring ABCA & line BCAB have a junction at B
func TestJoinRingLineReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("bcabLine", orb.LineString{
			orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0}, orb.Point{1, 0},
		}),
		NewTestFeature("abcaPolygon", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{1, 0}))
}

// join when an old arc ABC extends a new arc AB, there is a junction at B
func TestJoinOldArcExtends(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("ab", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a reversed old arc CBA extends a new arc AB, there is a junction at B
func TestJoinOldArcExtendsReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", orb.LineString{
			orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("ab", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a new arc ADE shares its start with an old arc ABC, there is a junction at A
func TestJoinNewArcSharesStart(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("ade", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 1}, orb.Point{2, 1},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{2, 1}))
}

// join ring ABA has no junctions
func TestJoinRingNoJunctions(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("aba", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join ring AA has no junctions
func TestJoinRingAANoJunctions(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("aa", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join degenerate ring A has no junctions
func TestJoinRingANoJunctions(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("a", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join when a new line DEC shares its end with an old line ABC, there is a junction at C
func TestJoinNewLineSharesEnd(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("dec", orb.LineString{
			orb.Point{0, 1}, orb.Point{1, 1}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{0, 1}))
}

// join when a new line ABC extends an old line AB, there is a junction at B
func TestJoinNewLineExtends(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("ab", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0},
		}),
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a new line ABC extends a reversed old line BA, there is a junction at B
func TestJoinNewLineExtendsReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("ba", orb.LineString{
			orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a new line starts BC in the middle of an old line ABC, there is a junction at B
func TestJoinNewStartsMiddle(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("bc", orb.LineString{
			orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a new line BC starts in the middle of a reversed old line CBA, there is a junction at B
func TestJoinNewStartsMiddleReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", orb.LineString{
			orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("bc", orb.LineString{
			orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a new line ABD deviates from an old line ABC, there is a junction at B
func TestJoinNewLineDeviates(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("abd", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{3, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 4)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{3, 0}))
}

// join when a new line ABD deviates from a reversed old line CBA, there is a junction at B
func TestJoinNewLineDeviatesReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", orb.LineString{
			orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("abd", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{3, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 4)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{3, 0}))
}

// join when a new line DBC merges into an old line ABC, there is a junction at B
func TestJoinNewLineMerges(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("dbc", orb.LineString{
			orb.Point{3, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 4)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{3, 0}))
}

// join when a new line DBC merges into a reversed old line CBA, there is a junction at B
func TestJoinNewLineMergesReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", orb.LineString{
			orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("dbc", orb.LineString{
			orb.Point{3, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 4)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{3, 0}))
}

// join when a new line DBE shares a single midpoint with an old line ABC, there is a junction at B
func TestJoinNewLineSharesMidpoint(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0},
		}),
		NewTestFeature("dbe", orb.LineString{
			orb.Point{0, 1}, orb.Point{1, 0}, orb.Point{2, 1},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 5)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{0, 1}))
	is.True(junctions.Has([]float64{2, 1}))
}

// join when a new line ABDE skips a point with an old line ABCDE, there is a junction at B and D
func TestJoinNewLineSkipsPoint(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcde", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{3, 0}, orb.Point{4, 0},
		}),
		NewTestFeature("adbe", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{3, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 4)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{3, 0}))
	is.True(junctions.Has([]float64{4, 0}))
}

// join when a new line ABDE skips a point with a reversed old line EDCBA, there is a junction at B and D
func TestJoinNewLineSkipsPointReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("edcba", orb.LineString{
			orb.Point{4, 0}, orb.Point{3, 0}, orb.Point{2, 0}, orb.Point{1, 0}, orb.Point{0, 0},
		}),
		NewTestFeature("adbe", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{3, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 4)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{3, 0}))
	is.True(junctions.Has([]float64{4, 0}))
}

// join when a line ABCDBE self-intersects with its middle, there are no junctions
func TestJoinSelfIntersectsMiddle(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbe", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{3, 0}, orb.Point{1, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 2)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{4, 0}))
}

// join when a line ABACD self-intersects with its start, there are no junctions
func TestJoinSelfIntersectsStart(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abacd", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 0}, orb.Point{3, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 2)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{4, 0}))
}

// join when a line ABCDBD self-intersects with its end, there are no junctions
func TestJoinSelfIntersectsEnd(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbd", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{4, 0}, orb.Point{3, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 2)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{4, 0}))
}

// join when an old line ABCDBE self-intersects and shares a point B, there is a junction at B
func TestJoinSelfIntersectsShares(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbe", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{3, 0}, orb.Point{1, 0}, orb.Point{4, 0},
		}),
		NewTestFeature("fbg", orb.LineString{
			orb.Point{0, 1}, orb.Point{1, 0}, orb.Point{2, 1},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 5)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{0, 1}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 1}))
	is.True(junctions.Has([]float64{4, 0}))
}

// join when a line ABCA is closed, there is a junction at A
func TestJoinLineClosed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{0, 0}))
}

// join when a ring ABCA is closed, there are no junctions
func TestJoinRingClosed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join exact duplicate rings ABCA & ABCA share the arc ABCA
func TestJoinDuplicateRingsShare(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join reversed duplicate rings ABCA & ACBA share the arc ABCA
func TestJoinDuplicateRingsReversedShare(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("acba", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{0, 1}, orb.Point{1, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join coincident rings ABCA & BCAB share the arc BCAB
func TestJoinCoincidentRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("bcab", orb.Polygon{
			orb.Ring{
				orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0}, orb.Point{1, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join coincident rings ABCA & BACB share the arc BCAB
func TestJoinCoincidentRings2(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("bacb", orb.Polygon{
			orb.Ring{
				orb.Point{1, 0}, orb.Point{0, 0}, orb.Point{0, 1}, orb.Point{1, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join coincident rings ABCA & DBED share the point B
func TestJoinCoincidentRingsShare(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("dbed", orb.Polygon{
			orb.Ring{
				orb.Point{2, 1}, orb.Point{1, 0}, orb.Point{2, 2}, orb.Point{2, 1},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{1, 0}))
}

// join coincident ring ABCA & line DBE share the point B
func TestJoinCoincidentRingLine(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("dbe", orb.LineString{
			orb.Point{2, 1}, orb.Point{1, 0}, orb.Point{2, 2},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 1}))
	is.True(junctions.Has([]float64{2, 2}))
}
