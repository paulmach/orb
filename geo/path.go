package geo

import (
	"bytes"
	"fmt"
	"math"

	"github.com/paulmach/go.geojson"
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

// NewPathFromEncoding is the inverse of path.Encode. It takes a string encoding of a lat/lon path
// and returns the actual path it represents. Factor defaults to 1.0e5,
// the same used by Google for polyline encoding.
func NewPathFromEncoding(encoded string, factor ...int) Path {
	var count, index int

	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	p := NewPath()
	tempLatLon := [2]int{0, 0}

	for index < len(encoded) {
		var result int
		var b = 0x20
		var shift uint

		for b >= 0x20 {
			b = int(encoded[index]) - 63
			index++

			result |= (b & 0x1f) << shift
			shift += 5
		}

		// sign dection
		if result&1 != 0 {
			result = ^(result >> 1)
		} else {
			result = result >> 1
		}

		if count%2 == 0 {
			result += tempLatLon[0]
			tempLatLon[0] = result
		} else {
			result += tempLatLon[1]
			tempLatLon[1] = result

			p = append(p, Point{float64(tempLatLon[1]) / f, float64(tempLatLon[0]) / f})
		}

		count++
	}

	return p
}

// NewPathFromXYData creates a path from a slice of [2]float64 values
// representing [horizontal, vertical] type data, for example lon/lat values from geojson.
func NewPathFromXYData(data [][2]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		p = append(p, Point{data[i][0], data[i][1]})
	}

	return p
}

// NewPathFromYXData creates a path from a slice of [2]float64 values
// representing [vertical, horizontal] type data, for example typical lat/lon data.
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

// Encode converts the path to a string using the Google Maps Polyline Encoding method.
// Factor defaults to 1.0e5, the same used by Google for polyline encoding.
func (p Path) Encode(factor ...int) string {
	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	var pLat int
	var pLon int

	var result bytes.Buffer
	scratch1 := make([]byte, 0, 50)
	scratch2 := make([]byte, 0, 50)

	for _, p := range p {
		lat5 := int(math.Floor(p.Lat()*f + 0.5))
		lon5 := int(math.Floor(p.Lon()*f + 0.5))

		deltaLat := lat5 - pLat
		deltaLon := lon5 - pLon

		pLat = lat5
		pLon = lon5

		result.Write(append(encodeSignedNumber(deltaLat, scratch1), encodeSignedNumber(deltaLon, scratch2)...))

		scratch1 = scratch1[:0]
		scratch2 = scratch2[:0]
	}

	return result.String()
}

func encodeSignedNumber(num int, result []byte) []byte {
	shiftedNum := num << 1

	if num < 0 {
		shiftedNum = ^shiftedNum
	}

	for shiftedNum >= 0x20 {
		result = append(result, byte(0x20|(shiftedNum&0x1f)+63))
		shiftedNum >>= 5
	}

	return append(result, byte(shiftedNum+63))
}

// Distance computes the total distance using spherical geometry.
func (p Path) Distance(haversine ...bool) float64 {
	yesgeo := yesHaversine(haversine)
	sum := 0.0

	loopTo := len(p) - 1
	for i := 0; i < loopTo; i++ {
		sum += p[i].DistanceFrom(p[i+1], yesgeo)
	}

	return sum
}

// Bound returns a bound around the path. Uses rectangular coordinates.
func (p Path) Bound() Bound {
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

// GeoJSON creates a new geojson feature with a linestring geometry
// containing all the points.
func (p Path) GeoJSON() *geojson.Feature {
	coords := make([][]float64, 0, len(p))

	for _, point := range p {
		coords = append(coords, []float64{point[0], point[1]})
	}

	return geojson.NewLineStringFeature(coords)
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

// String returns the wkt representation of the path.
func (p Path) String() string {
	return p.WKT()
}
