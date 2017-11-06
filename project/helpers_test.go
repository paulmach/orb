package project

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestProjectionTypes(t *testing.T) {
	ToPlanar(orb.Point{0, 0}, Mercator)
	ToGeo(orb.Point{0, 0}, Mercator)
}
