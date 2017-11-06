package project

import (
	"testing"

	"github.com/paulmach/orb/geo"
)

func TestProjectionTypes(t *testing.T) {
	ToPlanar(geo.Point{0, 0}, Mercator)
	ToGeo(geo.Point{0, 0}, Mercator)
}
