package planar

import (
	"bytes"
	"fmt"
	"math"
)

// Path represents a set of points to be thought of as a polyline.
type Path []Point

// NewPath creates a new path.
func NewPath() Path {
	return NewPathPreallocate(0, 100)
}

// NewPathPreallocate creates a new path with points array of the given size.
func NewPathPreallocate(length, capacity int) Path {
	return Path(make([]Point, length, capacity))
}

// NewPathFromXYData creates a path from a slice of [2]float64 values
// representing [horizontal, vertical] type data, for example lng/lat values from geojson.
func NewPathFromXYData(data [][2]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		p = append(p, Point{data[i][0], data[i][1]})
	}

	return p
}

// NewPathFromYXData creates a path from a slice of [2]float64 values
// representing [vertical, horizontal] type data, for example typical lat/lng data.
func NewPathFromYXData(data [][2]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		p = append(p, Point{data[i][1], data[i][0]})
	}

	return p
}

// NewPathFromXYSlice creates a path from a slice of []float64 values.
// The first two elements are taken to be horizontal and vertical components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewPathFromXYSlice(data [][]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			p = append(p, Point{data[i][0], data[i][1]})
		}
	}

	return p
}

// NewPathFromYXSlice creates a path from a slice of []float64 values.
// The first two elements are taken to be vertical and horizontal components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewPathFromYXSlice(data [][]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			p = append(p, Point{data[i][1], data[i][0]})
		}
	}

	return p
}

// Distance computes the total distance in the units of the points.
func (p Path) Distance() float64 {
	sum := 0.0

	loopTo := len(p) - 1
	for i := 0; i < loopTo; i++ {
		sum += p[i].DistanceFrom(p[i+1])
	}

	return sum
}

// DistanceFrom computes an O(n) distance from the path. Loops over every
// subline to find the minimum distance.
func (p Path) DistanceFrom(point Point) float64 {
	return math.Sqrt(p.DistanceFromSquared(point))
}

// DistanceFromSquared computes an O(n) minimum squared distance from the path.
// Loops over every subline to find the minimum distance.
func (p Path) DistanceFromSquared(point Point) float64 {
	dist := math.Inf(1)

	l := Line{}
	loopTo := len(p) - 1
	for i := 0; i < loopTo; i++ {
		l.a = p[i]
		l.b = p[i+1]
		dist = math.Min(l.DistanceFromSquared(point), dist)
	}

	return dist
}

// Interpolate interpolates the path by geo distance.
func (p Path) Interpolate(percent float64) Point {
	if percent <= 0 {
		return p[0]
	} else if percent >= 1 {
		return p[len(p)-1]
	}

	destination := p.Distance() * percent
	travelled := 0.0

	for i := 0; i < len(p)-1; i++ {
		dist := p[i].DistanceFrom(p[i+1])
		if (travelled + dist) > destination {
			factor := (destination - travelled) / dist
			return Point{
				p[i][0]*(1-factor) + p[i+1][0]*factor,
				p[i][1]*(1-factor) + p[i+1][1]*factor,
			}
		}
		travelled += dist
	}

	return p[0]
}

// Project computes the percent along this path closest to the given point,
// normalized to the length of the path.
func (p Path) Project(point Point) float64 {
	minDistance := math.Inf(1)
	measure := math.Inf(-1)
	sum := 0.0

	seg := Line{}
	for i := 0; i < len(p)-1; i++ {
		seg.a = p[i]
		seg.b = p[i+1]

		distanceToLine := seg.DistanceFromSquared(point)
		segDistance := seg.Distance()

		if distanceToLine < minDistance {
			minDistance = distanceToLine

			proj := seg.Project(point)
			if proj < 0 {
				proj = 0
			} else if proj > 1 {
				proj = 1
			}

			measure = sum + proj*segDistance
		}
		sum += segDistance
	}
	return measure / sum
}

// Bound returns a rectangle bound around the path. Uses rectangular coordinates.
func (p Path) Bound() Rect {
	return PointSet(p).Bound()
}

// Equal compares two paths. Returns true if lengths are the same
// and all points are Equal.
func (p Path) Equal(path Path) bool {
	return PointSet(p).Equal(PointSet(path))
}

// Clone returns a new copy of the path.
func (p Path) Clone() Path {
	ps := PointSet(p)
	return Path(ps.Clone())
}

// WKT returns the path in WKT format, eg. LINESTRING(30 10,10 30,40 40)
// For empty paths the result will be 'EMPTY'.
func (p Path) WKT() string {
	if len(p) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "LINESTRING(%g %g", p[0][0], p[0][1])

	for i := 1; i < len(p); i++ {
		fmt.Fprintf(buff, ",%g %g", p[i][0], p[i][1])
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// String returns a string representation of the path.
// The format is WKT, e.g. LINESTRING(30 10,10 30,40 40)
// For empty paths the result will be 'EMPTY'.
func (p Path) String() string {
	return p.WKT()
}
