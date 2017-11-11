package orb

import (
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	for _, g := range AllGeometries {
		// this closure is necessary if tests are run in parallel, maybe
		func(geom Geometry) {
			t.Run(fmt.Sprintf("%T", g), func(t *testing.T) {
				if !Equal(geom, geom) {
					t.Errorf("%T not equal", g)
				}
			})
		}(g)
	}
}
