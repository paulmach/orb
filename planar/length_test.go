package planar

import (
	"testing"

	"github.com/paulmach/orb/geo"
)

func TestLength_LineString(t *testing.T) {
	ls := geo.LineString{{0, 0}, {0, 3}, {4, 3}}
	if d := Length(ls); d != 7 {
		t.Errorf("length got: %f != 7.0", d)
	}
}
