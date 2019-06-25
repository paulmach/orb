package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	orb "github.com/paulmach/orb"
	geojson "github.com/paulmach/orb/geojson"
)

// Sets the quantization transform
func TestPreQuantize(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{}
	topo := &Topology{
		BBox:  []float64{0, 0, 1, 1},
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  1e4,
			PostQuantize: 1e4,
		},
	}

	topo.preQuantize()

	is.Equal(topo.Transform, &Transform{
		Scale:     [2]float64{float64(1) / 9999, float64(1) / 9999},
		Translate: [2]float64{0, 0},
	})
}

// Converts coordinates to fixed precision
func TestPreQuantizeConverts(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
		}),
	}

	expected := orb.LineString{
		orb.Point{0, 0}, orb.Point{9999, 0}, orb.Point{0, 9999}, orb.Point{0, 0},
	}

	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  1e4,
			PostQuantize: 1e4,
		},
	}

	topo.bounds()
	topo.preQuantize()

	is.Equal(topo.Transform, &Transform{
		Scale:     [2]float64{float64(1) / 9999, float64(1) / 9999},
		Translate: [2]float64{0, 0},
	})
	is.Equal(topo.input[0].Geometry, expected)
}

// Observes the quantization parameter
func TestPreQuantizeObserves(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
		}),
	}

	expected := orb.LineString{
		{0, 0}, {9, 0}, {0, 9}, {0, 0},
	}

	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  10,
			PostQuantize: 10,
		},
	}

	topo.bounds()
	topo.preQuantize()
	is.Equal(topo.input[0].Geometry, expected)
}

// Observes the bounding box
func TestPreQuantizeObservesBB(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{1, 0}, orb.Point{0, 1}, orb.Point{0, 0},
		}),
	}
	topo := &Topology{
		BBox:  []float64{-1, -1, 2, 2},
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  10,
			PostQuantize: 10,
		},
	}

	topo.preQuantize()

	expected := orb.LineString{
		{3, 3}, {6, 3}, {3, 6}, {3, 3},
	}
	is.Equal(topo.input[0].Geometry, expected)
}

// Applies to points as well as arcs
func TestPreQuantizeAppliesToPoints(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.MultiPoint{
			{0, 0}, {1, 0}, {0, 1}, {0, 0},
		})}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  1e4,
			PostQuantize: 1e4,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := orb.MultiPoint{
		{0, 0}, {9999, 0}, {0, 9999}, {0, 0},
	}
	is.Equal(topo.input[0].Geometry, expected)
}

// Skips coincident points in lines
func TestPreQuantizeSkipsCoincidencesInLines(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{0, 0}, orb.Point{0.9, 0.9}, orb.Point{1.1, 1.1}, orb.Point{2, 2},
		}),
	}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := orb.LineString{
		{0, 0}, {1, 1}, {2, 2},
	}
	is.Equal(topo.input[0].Geometry, expected)
}

// Skips coincident points in polygons
func TestPreQuantizeSkipsCoincidencesInPolygons(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("polygon", orb.Polygon{
			orb.Ring{
				orb.Point{0, 0}, orb.Point{0.9, 0.9}, orb.Point{1.1, 1.1}, orb.Point{2, 2}, orb.Point{0, 0},
			},
		}),
	}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := orb.Polygon{
		orb.Ring{
			{0, 0}, {1, 1}, {2, 2}, {0, 0},
		},
	}

	is.Equal(topo.input[0].Geometry, expected)
}

// Does not skip coincident points in points
func TestPreQuantizeDoesntSkipInPoints(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("multipoint", orb.MultiPoint{
			{0, 0}, {0.9, 0.9}, {1.1, 1.1}, {2, 2}, {0, 0},
		}),
	}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := orb.MultiPoint{
		{0, 0}, {1, 1}, {1, 1}, {2, 2}, {0, 0},
	}

	is.Equal(topo.input[0].Geometry, expected)
}

// Includes closing point in degenerate lines
func TestPreQuantizeIncludesClosingLine(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", orb.LineString{
			orb.Point{1, 1}, orb.Point{1, 1}, orb.Point{1, 1},
		}),
	}
	topo := &Topology{
		BBox:  []float64{0, 0, 2, 2},
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.preQuantize()

	expected := []orb.Point{
		{1, 1}, {1, 1},
	}

	is.Equal(topo.input[0].Geometry, expected)
}

// Includes closing point in degenerate polygons
func TestPreQuantizeIncludesClosingPolygon(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("polygon", orb.Polygon{
			orb.Ring{
				orb.Point{0.9, 1}, orb.Point{1.1, 1}, orb.Point{1.01, 1}, orb.Point{0.9, 1},
			},
		}),
	}
	topo := &Topology{
		BBox:  []float64{0, 0, 2, 2},
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.preQuantize()

	expected := orb.Polygon{
		orb.Ring{
			{1, 1}, {1, 1},
		},
	}

	is.Equal(topo.input[0].Geometry, expected)
}
