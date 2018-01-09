package wkb

import (
	"encoding/binary"
	"testing"

	"github.com/paulmach/orb"
)

func TestCentroidArea(t *testing.T) {
	for _, g := range orb.AllGeometries {
		Marshal(g, binary.BigEndian)
	}
}
