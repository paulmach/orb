package project

import (
	"testing"

	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

func TestProjectionTypes(t *testing.T) {
	ToPlanar(geo.Point{0, 0}, Mercator)
	ToGeo(planar.Point{0, 0}, Mercator)
}
