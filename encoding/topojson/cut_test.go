package topojson

import (
	"reflect"
	"testing"

	"github.com/cheekybits/is"
	orb "github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// See https://github.com/mbostock/topojson/blob/master/test/topology/cut-test.js

// cut exact duplicate lines ABC & ABC have no cuts
func TestCutDuplicates(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc2")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 5)
	is.Nil(o2.Arc.Next)
}

// cut reversed duplicate lines ABC & CBA have no cuts
func TestCutReversedDuplicates(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "cba")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 5)
	is.Nil(o2.Arc.Next)
}

// cut exact duplicate rings ABCA & ABCA have no cuts
func TestCutDuplicateRings(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "abca2")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut reversed duplicate rings ACBA & ABCA have no cuts
func TestCutReversedDuplicateRings(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "acba")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut rotated duplicate rings BCAB & ABCA have no cuts
func TestCutRotatedDuplicateRings(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "bcab")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut ring ABCA & line ABCA have no cuts
func TestCutRingLine(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abcaLine")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abcaPolygon")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut ring BCAB & line ABCA have no cuts
func TestCutRingLineReversed(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abcaLine")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "bcabPolygon")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut ring ABCA & line BCAB have no cuts
func TestCutRingLineReversed2(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("bcabLine", orb.LineString{
			orb.Point{1, 0}, {2, 0}, orb.Point{0, 0}, orb.Point{1, 0},
		}),
		NewTestFeature("abcaPolygon", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "bcabLine")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abcaPolygon")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut when an old arc ABC extends a new arc AB, ABC is cut into AB-BC
func TestCutOldArcExtends(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "ab")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Nil(o2.Arc.Next)
}

// cut when a reversed old arc CBA extends a new arc AB, CBA is cut into CB-BA
func TestCutReversedOldArcExtends(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "ab")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Nil(o2.Arc.Next)
}

// cut when a new arc ADE shares its start with an old arc ABC, there are no cuts
func TestCutNewArcSharesStart(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "ade")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 5)
	is.Nil(o2.Arc.Next)
}

// cut ring ABA has no junctions
func TestCutRingNoCuts(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "aba")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 2)
	is.Nil(o1.Arcs[0].Next)
}

// cut ring AA has no cuts
func TestCutRingAANoCuts(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "aa")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 1)
	is.Nil(o1.Arcs[0].Next)
}

// cut degenerate ring A has no cuts
func TestCutRingANoCuts(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "a")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 0)
	is.Nil(o1.Arcs[0].Next)
}

// cut when a new line DEC shares its end with an old line ABC, there are no cuts
func TestCutNewLineSharesEnd(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "dec")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 5)
	is.Nil(o2.Arc.Next)
}

// cut when a new line ABC extends an old line AB, ABC is cut into AB-BC
func TestCutNewLineExtends(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "ab")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 2)
	is.Equal(o2.Arc.End, 3)
	is.Equal(o2.Arc.Next.Start, 3)
	is.Equal(o2.Arc.Next.End, 4)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a new line ABC extends a reversed old line BA, ABC is cut into AB-BC
func TestCutNewLineExtendsReversed(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "ba")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 2)
	is.Equal(o2.Arc.End, 3)
	is.Equal(o2.Arc.Next.Start, 3)
	is.Equal(o2.Arc.Next.End, 4)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a new line starts BC in the middle of an old line ABC, ABC is cut into AB-BC
func TestCutNewStartsMiddle(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "bc")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Nil(o2.Arc.Next)
}

// cut when a new line BC starts in the middle of a reversed old line CBA, CBA is cut into CB-BA
func TestCutNewStartsMiddleReversed(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "bc")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Nil(o2.Arc.Next)
}

// cut when a new line ABD deviates from an old line ABC, ABD is cut into AB-BD and ABC is cut into AB-BC
func TestCutNewLineDeviates(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "abd")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a new line ABD deviates from a reversed old line CBA, CBA is cut into CB-BA and ABD is cut into AB-BD
func TestCutNewLineDeviatesReversed(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "abd")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a new line DBC merges into an old line ABC, DBC is cut into DB-BC and ABC is cut into AB-BC
func TestCutNewLineMerges(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "dbc")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a new line DBC merges into a reversed old line CBA, DBC is cut into DB-BC and CBA is cut into CB-BA
func TestCutNewLineMergesReversed(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "dbc")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a new line DBE shares a single midpoint with an old line ABC, DBE is cut into DB-BE and ABC is cut into AB-BC
func TestCutNewLineSharesMidpoint(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "dbe")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a new line ABDE skips a point with an old line ABCDE, ABDE is cut into AB-BD-DE and ABCDE is cut into AB-BCD-DE
func TestCutNewLineSkipsPoint(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abcde")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	o1Next := o1.Arc.Next
	is.Equal(o1Next.Start, 1)
	is.Equal(o1Next.End, 3)
	is.Equal(o1Next.Next.Start, 3)
	is.Equal(o1Next.Next.End, 4)
	is.Nil(o1Next.Next.Next)

	o2 := GetFeature(topo, "adbe")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 5)
	is.Equal(o2.Arc.End, 6)
	o2Next := o2.Arc.Next
	is.Equal(o2Next.Start, 6)
	is.Equal(o2Next.End, 7)
	is.Equal(o2Next.Next.Start, 7)
	is.Equal(o2Next.Next.End, 8)
	is.Nil(o2Next.Next.Next)
}

// cut when a new line ABDE skips a point with a reversed old line EDCBA, ABDE is cut into AB-BD-DE and EDCBA is cut into ED-DCB-BA
func TestCutNewLineSkipsPointReversed(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "edcba")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	o1Next := o1.Arc.Next
	is.Equal(o1Next.Start, 1)
	is.Equal(o1Next.End, 3)
	is.Equal(o1Next.Next.Start, 3)
	is.Equal(o1Next.Next.End, 4)
	is.Nil(o1Next.Next.Next)

	o2 := GetFeature(topo, "adbe")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 5)
	is.Equal(o2.Arc.End, 6)
	o2Next := o2.Arc.Next
	is.Equal(o2Next.Start, 6)
	is.Equal(o2Next.End, 7)
	is.Equal(o2Next.Next.Start, 7)
	is.Equal(o2Next.Next.End, 8)
	is.Nil(o2Next.Next.Next)
}

// cut when a line ABCDBE self-intersects with its middle, it is not cut
func TestCutSelfIntersectsMiddle(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbe", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{3, 0}, orb.Point{1, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "abcdbe")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 5)
	is.Nil(o1.Arc.Next)
}

// cut when a line ABACD self-intersects with its start, it is cut into ABA-ACD
func TestCutSelfIntersectsStart(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abacd", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 0}, orb.Point{3, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "abacd")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Equal(o1.Arc.Next.Start, 2)
	is.Equal(o1.Arc.Next.End, 4)
	is.Nil(o1.Arc.Next.Next)
}

// cut when a line ABDCD self-intersects with its end, it is cut into ABD-DCD
func TestCutSelfIntersectsEnd(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbd", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{4, 0}, orb.Point{3, 0}, orb.Point{4, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "abcdbd")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Equal(o1.Arc.Next.Start, 2)
	is.Equal(o1.Arc.Next.End, 4)
	is.Nil(o1.Arc.Next.Next)
}

// cut when an old line ABCDBE self-intersects and shares a point B, ABCDBE is cut into AB-BCDB-BE and FBG is cut into FB-BG
func TestCutSelfIntersectsShares(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abcdbe")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	o1Next := o1.Arc.Next
	is.Equal(o1Next.Start, 1)
	is.Equal(o1Next.End, 4)
	is.Equal(o1Next.Next.Start, 4)
	is.Equal(o1Next.Next.End, 5)
	is.Nil(o1Next.Next.Next)

	o2 := GetFeature(topo, "fbg")
	is.Equal(o2.Type, geojson.TypeLineString)
	is.Equal(o2.Arc.Start, 6)
	is.Equal(o2.Arc.End, 7)
	is.Equal(o2.Arc.Next.Start, 7)
	is.Equal(o2.Arc.Next.End, 8)
	is.Nil(o2.Arc.Next.Next)
}

// cut when a line ABCA is closed, there are no cuts
func TestCutLineClosed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypeLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)
}

// cut when a ring ABCA is closed, there are no cuts
func TestCutRingClosed(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)
}

// cut exact duplicate rings ABCA & ABCA have no cuts
func TestCutDuplicateRingsShare(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("abca2", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "abca2")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut reversed duplicate rings ABCA & ACBA have no cuts
func TestCutDuplicateRingsReversedShare(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "acba")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut coincident rings ABCA & BCAB have no cuts
func TestCutCoincidentRings(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "bcab")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut coincident rings ABCA & BACB have no cuts
func TestCutCoincidentRings2(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "bacb")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)
}

// cut coincident rings ABCDA, EFAE & GHCG are cut into ABC-CDA, EFAE and GHCG
func TestCutCoincidentRings3(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcda", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{1, 1}, orb.Point{0, 1}, orb.Point{0, 0},
			},
		}),
		NewTestFeature("efae", orb.Polygon{
			orb.Ring{
				orb.Point{0, -1}, orb.Point{1, -1}, orb.Point{0, 0}, orb.Point{0, -1},
			},
		}),
		NewTestFeature("ghcg", orb.Polygon{
			orb.Ring{
				orb.Point{0, 2}, orb.Point{1, 2}, orb.Point{1, 1}, orb.Point{0, 2},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "abcda")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 2)
	is.Equal(o1.Arcs[0].Next.Start, 2)
	is.Equal(o1.Arcs[0].Next.End, 4)
	is.Nil(o1.Arcs[0].Next.Next)

	o2 := GetFeature(topo, "efae")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 5)
	is.Equal(o2.Arcs[0].End, 8)
	is.Nil(o2.Arcs[0].Next)

	o3 := GetFeature(topo, "ghcg")
	is.Equal(o3.Type, geojson.TypePolygon)
	is.Equal(len(o3.Arcs), 1)
	is.Equal(o3.Arcs[0].Start, 9)
	is.Equal(o3.Arcs[0].End, 12)
	is.Nil(o3.Arcs[0].Next)
}

// cut coincident rings ABCA & DBED have no cuts, but are rotated to share B
func TestCutNoCutsButRotated(t *testing.T) {
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
	topo.cut()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "dbed")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)

	is.True(reflect.DeepEqual(topo.coordinates[0:4], [][]float64{
		{1, 0}, {0, 1}, {0, 0}, {1, 0},
	}))
	is.True(reflect.DeepEqual(topo.coordinates[4:8], [][]float64{
		{1, 0}, {2, 2}, {2, 1}, {1, 0},
	}))
}

// cut overlapping rings ABCDA and BEFCB are cut into BC-CDAB and BEFC-CB
func TestCutOverlapping(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcda", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{1, 1}, orb.Point{0, 1}, orb.Point{0, 0}, // rotated to BCDAB, cut BC-CDAB
			},
		}),
		NewTestFeature("befcb", orb.Polygon{
			orb.Ring{
				orb.Point{1, 0}, orb.Point{2, 0}, orb.Point{2, 1}, orb.Point{1, 1}, orb.Point{1, 0},
			},
		}),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()

	o1 := GetFeature(topo, "abcda")
	is.Equal(o1.Type, geojson.TypePolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 1)
	is.Equal(o1.Arcs[0].Next.Start, 1)
	is.Equal(o1.Arcs[0].Next.End, 4)
	is.Nil(o1.Arcs[0].Next.Next)

	o2 := GetFeature(topo, "befcb")
	is.Equal(o2.Type, geojson.TypePolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 5)
	is.Equal(o2.Arcs[0].End, 8)
	is.Equal(o2.Arcs[0].Next.Start, 8)
	is.Equal(o2.Arcs[0].Next.End, 9)
	is.Nil(o2.Arcs[0].Next.Next)
}
