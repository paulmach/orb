package geo

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/paulmach/go.geojson"
)

// LineString represents a set of points to be thought of as a polyline.
type LineString []Point

// NewLineString creates a new line string.
func NewLineString() LineString {
	return NewLineStringPreallocate(0, 100)
}

// NewLineStringPreallocate creates a new line string with points array of the given size.
func NewLineStringPreallocate(length, capacity int) LineString {
	return LineString(make([]Point, length, capacity))
}

// NewLineStringFromEncoding is the inverse of lineString.Encode. It takes a string encoding
// of a lat/lon path and returns the actual path it represents. Factor defaults to 1.0e5,
// the same used by Google for polyline encoding.
func NewLineStringFromEncoding(encoded string, factor ...int) LineString {
	var count, index int

	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	ls := NewLineString()
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

			ls = append(ls, Point{float64(tempLatLon[1]) / f, float64(tempLatLon[0]) / f})
		}

		count++
	}

	return ls
}

// NewLineStringFromXYData creates a line string from a slice of [2]float64 values
// representing [horizontal, vertical] type data, for example lon/lat values from geojson.
func NewLineStringFromXYData(data [][2]float64) LineString {
	ls := NewLineStringPreallocate(0, len(data))
	for i := range data {
		ls = append(ls, Point{data[i][0], data[i][1]})
	}

	return ls
}

// NewLineStringFromYXData creates a line string from a slice of [2]float64 values
// representing [vertical, horizontal] type data, for example typical lat/lon data.
func NewLineStringFromYXData(data [][2]float64) LineString {
	ls := NewLineStringPreallocate(0, len(data))
	for i := range data {
		ls = append(ls, Point{data[i][1], data[i][0]})
	}

	return ls
}

// NewLineStringFromXYSlice creates a line string from a slice of []float64 values.
// The first two elements are taken to be horizontal and vertical components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewLineStringFromXYSlice(data [][]float64) LineString {
	ls := NewLineStringPreallocate(0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			ls = append(ls, Point{data[i][0], data[i][1]})
		}
	}

	return ls
}

// NewLineStringFromYXSlice creates a line string from a slice of []float64 values.
// The first two elements are taken to be vertical and horizontal components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewLineStringFromYXSlice(data [][]float64) LineString {
	ls := NewLineStringPreallocate(0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			ls = append(ls, Point{data[i][1], data[i][0]})
		}
	}

	return ls
}

// GeoJSONType returns the GeoJSON type for the object.
func (ls LineString) GeoJSONType() string {
	return "LineString"
}

// Encode converts the line string to a string using the Google Maps Polyline Encoding method.
// Factor defaults to 1.0e5, the same used by Google for polyline encoding.
func (ls LineString) Encode(factor ...int) string {
	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	var pLat int
	var pLon int

	var result bytes.Buffer
	scratch1 := make([]byte, 0, 50)
	scratch2 := make([]byte, 0, 50)

	for _, ls := range ls {
		lat5 := int(math.Floor(ls.Lat()*f + 0.5))
		lon5 := int(math.Floor(ls.Lon()*f + 0.5))

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
func (ls LineString) Distance(haversine ...bool) float64 {
	yesgeo := yesHaversine(haversine)
	sum := 0.0

	loopTo := len(ls) - 1
	for i := 0; i < loopTo; i++ {
		sum += ls[i].DistanceFrom(ls[i+1], yesgeo)
	}

	return sum
}

// Reverse changes the direction of the line string.
// It returns a new list string.
func (ls LineString) Reverse() LineString {
	n := NewLineStringPreallocate(len(ls), len(ls))

	l := len(n) - 1
	for i := 0; i <= l/2; i++ {
		n[i], n[l-i] = ls[l-i], ls[i]
	}

	return n
}

// Bound returns a rect around the line string. Uses rectangular coordinates.
func (ls LineString) Bound() Rect {
	return MultiPoint(ls).Bound()
}

// Equal compares two line strings. Returns true if lengths are the same
// and all points are Equal.
func (ls LineString) Equal(lineString LineString) bool {
	return MultiPoint(ls).Equal(MultiPoint(lineString))
}

// Clone returns a new copy of the line string.
func (ls LineString) Clone() LineString {
	ps := MultiPoint(ls)
	return LineString(ps.Clone())
}

// GeoJSON creates a new geojson feature with a linestring geometry
// containing all the points.
func (ls LineString) GeoJSON() *geojson.Feature {
	coords := make([][]float64, 0, len(ls))

	for _, point := range ls {
		coords = append(coords, []float64{point[0], point[1]})
	}

	return geojson.NewLineStringFeature(coords)
}

// WKT returns the line string in WKT format, eg. LINESTRING(30 10,10 30,40 40)
// For empty line strings the result will be 'EMPTY'.
func (ls LineString) WKT() string {
	if len(ls) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "LINESTRING")
	wktPoints(buff, ls)

	return buff.String()
}

// String returns the wkt representation of the line string.
func (ls LineString) String() string {
	return ls.WKT()
}

func wktPoints(w io.Writer, ps []Point) {
	fmt.Fprintf(w, "(%g %g", ps[0][0], ps[0][1])

	for i := 1; i < len(ps); i++ {
		fmt.Fprintf(w, ",%g %g", ps[i][0], ps[i][1])
	}

	fmt.Fprintf(w, ")")
}
