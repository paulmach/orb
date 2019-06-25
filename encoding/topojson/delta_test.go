package topojson

import (
	"testing"

	"github.com/cheekybits/is"
)

// Converts arcs to delta encoding
func TestDeltaConverts(t *testing.T) {
	is := is.New(t)

	topo := &Topology{
		Arcs: [][][]float64{
			{
				{0, 0}, {9999, 0}, {0, 9999}, {0, 0},
			},
		},
		opts: &TopologyOptions{PostQuantize: 1e4},
	}

	expected := [][][]float64{
		{
			{0, 0}, {9999, 0}, {-9999, 9999}, {0, -9999},
		},
	}

	topo.delta()
	is.Equal(topo.Arcs, expected)
}

// Does not skip coincident points
func TestDeltaDoesntSkip(t *testing.T) {
	is := is.New(t)

	topo := &Topology{
		Arcs: [][][]float64{
			{
				{0, 0}, {9999, 0}, {9999, 0}, {0, 9999}, {0, 0},
			},
		},
		opts: &TopologyOptions{PostQuantize: 1e4},
	}

	expected := [][][]float64{
		{
			{0, 0}, {9999, 0}, {0, 0}, {-9999, 9999}, {0, -9999},
		},
	}

	topo.delta()
	is.Equal(topo.Arcs, expected)
}
