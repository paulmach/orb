package wkb

import "github.com/paulmach/orb"

// PointTestCase is used to test point scanning.
type PointTestCase struct {
	X, Y float64
	Data interface{}
	Err  error
}

// PointTestCases check different point scanning behavior.
var PointTestCases = []PointTestCase{
	{ // little endian
		X: -122.4546440212, Y: 37.7382859071,
		Data: []byte{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
	},
	{ // big endian
		X: -122.4546440212, Y: 37.7382859071,
		Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39, 17, 11, 205},
	},
	{ // mysql srid+wkb
		X: -122.671129, Y: 38.177484,
		Data: []byte{215, 15, 0, 0, 1, 1, 0, 0, 0, 107, 153, 12, 199, 243, 170, 94, 192, 25, 200, 179, 203, 183, 22, 67, 64},
	},
	{ // mysql srid+wkb, empty srid
		X: -122.671129, Y: 38.177484,
		Data: []byte{0, 0, 0, 0, 1, 1, 0, 0, 0, 107, 153, 12, 199, 243, 170, 94, 192, 25, 200, 179, 203, 183, 22, 67, 64},
	},
	{
		X: -93.787988, Y: 32.392335,
		Data: []byte{1, 1, 0, 0, 0, 253, 104, 56, 101, 110, 114, 87, 192, 192, 9, 133, 8, 56, 50, 64, 64},
	},
	{
		Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
		Err:  orb.ErrIncorrectGeometry,
	},
	{
		Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
		Err:  orb.ErrNotWKB,
	},
	{
		Data: []byte{0, 2, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
		Err:  orb.ErrIncorrectGeometry,
	},
	{
		Data: 123,
		Err:  orb.ErrUnsupportedDataType,
	},
}

// SegmentTestCase is used to test segment scanning.
type SegmentTestCase struct {
	Points [][2]float64
	Data   interface{}
	Err    error
}

// SegmentTestCases check different segment scanning behavior.
var SegmentTestCases = []SegmentTestCase{
	{
		Points: [][2]float64{{-123.016508, 38.040608}, {-122.670176, 38.548019}},
		Data:   []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
	},
	{
		Points: [][2]float64{{-123.016508, 38.040608}, {-122.670176, 38.548019}},
		Data:   []byte{215, 15, 0, 0, 1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
	},
	{
		Points: [][2]float64{{-72.796408, -45.407131}, {-72.688541, -45.384987}},
		Data:   []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 117, 145, 66, 89, 248, 50, 82, 192, 9, 24, 93, 222, 28, 180, 70, 192, 33, 61, 69, 14, 17, 44, 82, 192, 77, 49, 7, 65, 71, 177, 70, 192},
	},
	{
		Data: 123,
		Err:  orb.ErrUnsupportedDataType,
	},
	{
		Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
		Err:  orb.ErrIncorrectGeometry,
	},
	{
		Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
		Err:  orb.ErrNotWKB,
	},
	{
		Data: []byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
		Err:  orb.ErrIncorrectGeometry,
	},
}

// MultiPointTestCase is used to get point set scanning.
type MultiPointTestCase struct {
	Points [][2]float64
	Data   interface{}
	Err    error
}

// MultiPointTestCases check different point set scanning behavior.
var MultiPointTestCases = []MultiPointTestCase{
	{
		Points: [][2]float64{{1, 2}, {3, 4}, {5, 6}},
		Data:   []byte{0, 0, 0, 0, 1, 4, 0, 0, 0, 3, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 16, 64, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 64, 0, 0, 0, 0, 0, 0, 24, 64},
	},
	{
		Data: 123,
		Err:  orb.ErrUnsupportedDataType,
	},
	{
		Data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
		Err:  orb.ErrIncorrectGeometry,
	},
	{
		Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
		Err:  orb.ErrNotWKB,
	},
	{
		Data: []byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
		Err:  orb.ErrIncorrectGeometry,
	},
}

var testLineStringWKB = []byte{1, 2, 0, 0, 0, 6, 0, 0, 0, 205, 228, 155, 109, 110, 114, 87, 192, 174, 158, 147, 222, 55, 50, 64, 64, 134, 56, 214, 197, 109, 114, 87, 192, 238, 235, 192, 57, 35, 50, 64, 64, 173, 47, 18, 218, 114, 114, 87, 192, 25, 4, 86, 14, 45, 50, 64, 64, 10, 75, 60, 160, 108, 114, 87, 192, 224, 161, 40, 208, 39, 50, 64, 64, 149, 159, 84, 251, 116, 114, 87, 192, 96, 147, 53, 234, 33, 50, 64, 64, 195, 158, 118, 248, 107, 114, 87, 192, 89, 139, 79, 1, 48, 50, 64, 64}

// LineStringTestCase is used to get line string scanning.
type LineStringTestCase struct {
	Points [][2]float64
	Data   interface{}
	Err    error
}

// LineStringTestCases check different line string scanning behavior.
var LineStringTestCases = []LineStringTestCase{
	{
		Points: [][2]float64{{-93.78799, 32.39233}, {-93.78795, 32.3917}, {-93.78826, 32.392}, {-93.78788, 32.39184}, {-93.78839, 32.39166}, {-93.78784, 32.39209}},
		Data:   testLineStringWKB,
	},
	{
		Points: [][2]float64{{-93.78799, 32.39233}, {-93.78795, 32.3917}, {-93.78826, 32.392}, {-93.78788, 32.39184}, {-93.78839, 32.39166}, {-93.78784, 32.39209}},
		Data:   append([]byte{215, 15, 0, 0}, testLineStringWKB...),
	},
	{
		Points: [][2]float64{{-93.78799, 32.39233}, {-93.78795, 32.3917}, {-93.78826, 32.392}, {-93.78788, 32.39184}, {-93.78839, 32.39166}, {-93.78784, 32.39209}},
		Data:   append([]byte{0, 0, 0, 0}, testLineStringWKB...),
	},
	{
		Data: 123,
		Err:  orb.ErrUnsupportedDataType,
	},
	{
		Data: []byte{0, 0, 0, 0, 1},
		Err:  orb.ErrNotWKB,
	},
	{
		Data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
		Err:  orb.ErrNotWKB,
	},
	{
		Data: []byte{1, 1, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64},
		Err:  orb.ErrIncorrectGeometry,
	},
}
