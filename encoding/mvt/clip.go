package mvt

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/clip"
)

// Clip will clip all geometries in all layers to the given bounds.
func (ls Layers) Clip(box orb.Bound) {
	for _, l := range ls {
		l.Clip(box)
	}
}

// Clip will clip all geometries in this layer to the given bounds.
func (l *Layer) Clip(box orb.Bound) {
	for _, f := range l.Features {
		g := clip.Geometry(box, f.Geometry)
		f.Geometry = g
	}
}
