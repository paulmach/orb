package projections

import (
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
)

func ProjectPath(path geo.Path, f Project) planar.Path {
	n := planar.NewPathPreallocate(len(path), len(path))
	for i := range path {
		n[i] = f(path[i])
	}

	return n
}

func InvertPath(path planar.Path, f Inverse) geo.Path {
	n := geo.NewPathPreallocate(len(path), len(path))
	for i := range path {
		n[i] = f(path[i])
	}

	return n
}

func ProjectBound(bound geo.Bound, f Project) planar.Bound {
	return planar.NewBoundFromPoints(
		f(geo.Point(bound.SW)),
		f(geo.Point(bound.NE)),
	)

}

func InvertBound(bound planar.Bound, f Inverse) geo.Bound {
	return geo.NewBoundFromPoints(
		f(planar.Point(bound.SW)),
		f(planar.Point(bound.NE)),
	)
}
