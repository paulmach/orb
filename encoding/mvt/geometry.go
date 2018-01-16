package mvt

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt/vectortile"
	"github.com/pkg/errors"
)

const (
	moveTo    = 1
	lineTo    = 2
	closePath = 7
)

func encodeGeometry(g orb.Geometry) (vectortile.Tile_GeomType, []uint32, error) {
	switch g := g.(type) {
	case orb.Point:
		e := newGeomEncoder(3)
		e.MoveTo([]orb.Point{g})

		return vectortile.Tile_POINT, e.Data, nil
	case orb.MultiPoint:
		e := newGeomEncoder(1 + 2*len(g))
		e.MoveTo([]orb.Point(g))

		return vectortile.Tile_POINT, e.Data, nil
	case orb.LineString:
		e := newGeomEncoder(2 + 2*len(g))
		e.MoveTo([]orb.Point{g[0]})
		e.LineTo([]orb.Point(g[1:]))

		return vectortile.Tile_LINESTRING, e.Data, nil
	case orb.MultiLineString:
		e := newGeomEncoder(elMLS(g))
		for _, ls := range g {
			e.MoveTo([]orb.Point{ls[0]})
			e.LineTo([]orb.Point(ls[1:]))
		}

		return vectortile.Tile_LINESTRING, e.Data, nil
	case orb.Ring:
		e := newGeomEncoder(3 + 2*len(g))
		e.MoveTo([]orb.Point{g[0]})
		if g.Closed() {
			e.LineTo([]orb.Point(g[1 : len(g)-1]))
		} else {
			e.LineTo([]orb.Point(g[1:]))
		}
		e.ClosePath()

		return vectortile.Tile_POLYGON, e.Data, nil
	case orb.Polygon:
		e := newGeomEncoder(elP(g))
		for _, r := range g {
			e.MoveTo([]orb.Point{r[0]})
			if r.Closed() {
				e.LineTo([]orb.Point(r[1 : len(r)-1]))
			} else {
				e.LineTo([]orb.Point(r[1:]))
			}
			e.ClosePath()
		}

		return vectortile.Tile_POLYGON, e.Data, nil
	case orb.MultiPolygon:
		e := newGeomEncoder(elMP(g))
		for _, p := range g {
			for _, r := range p {
				e.MoveTo([]orb.Point{r[0]})
				if r.Closed() {
					e.LineTo([]orb.Point(r[1 : len(r)-1]))
				} else {
					e.LineTo([]orb.Point(r[1:]))
				}
				e.ClosePath()
			}
		}

		return vectortile.Tile_POLYGON, e.Data, nil
	case orb.Collection:
		return 0, nil, errors.New("geometry collections are not supported")
	case orb.Bound:
		return encodeGeometry(g.ToPolygon())
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

type geomEncoder struct {
	prevX, prevY int32
	Data         []uint32
}

func newGeomEncoder(l int) *geomEncoder {
	return &geomEncoder{
		Data: make([]uint32, 0, l),
	}
}

func (ge *geomEncoder) MoveTo(points []orb.Point) {
	l := uint32(len(points))
	ge.Data = append(ge.Data, (l<<3)|moveTo)
	ge.addPoints(points)
}

func (ge *geomEncoder) LineTo(points []orb.Point) {
	l := uint32(len(points))
	ge.Data = append(ge.Data, (l<<3)|lineTo)
	ge.addPoints(points)
}

func (ge *geomEncoder) addPoints(points []orb.Point) {
	for i := range points {
		x := int32(points[i][0]) - ge.prevX
		y := int32(points[i][1]) - ge.prevY

		ge.prevX = int32(points[i][0])
		ge.prevY = int32(points[i][1])

		ge.Data = append(ge.Data,
			uint32((x<<1)^(x>>31)),
			uint32((y<<1)^(y>>31)),
		)
	}
}

func (ge *geomEncoder) ClosePath() {
	ge.Data = append(ge.Data, (1<<3)|closePath)
}

type keyValueEncoder struct {
	Keys   []string
	keyMap map[string]uint32

	Values   []*vectortile.Tile_Value
	valueMap map[interface{}]uint32
}

// A geomDecoder holds state for geometry decoding.
type geomDecoder struct {
	geom []uint32
	i    int

	prev orb.Point
}

func decodeGeometry(geomType vectortile.Tile_GeomType, geom []uint32) (orb.Geometry, error) {
	if len(geom) < 2 {
		return nil, errors.Errorf("geom is not long enough: %v", geom)
	}

	gd := &geomDecoder{geom: geom}

	switch geomType {
	case vectortile.Tile_POINT:
		return gd.decodePoint()
	case vectortile.Tile_LINESTRING:
		return gd.decodeLineString()
	case vectortile.Tile_POLYGON:
		return gd.decodePolygon()
	}

	return nil, errors.Errorf("unknown geometry type: %v", geomType)
}

func (gd *geomDecoder) decodePoint() (orb.Geometry, error) {
	_, count, err := gd.cmdAndCount()
	if err != nil {
		return nil, err
	}

	if count == 1 {
		return gd.NextPoint(), nil
	}

	mp := make(orb.MultiPoint, 0, count)
	for i := uint32(0); i < count; i++ {
		mp = append(mp, gd.NextPoint())
	}

	return mp, nil
}

func (gd *geomDecoder) decodeLine() (orb.LineString, error) {
	cmd, count, err := gd.cmdAndCount()
	if err != nil {
		return nil, err
	}

	if cmd != moveTo || count != 1 {
		return nil, errors.New("first command not one moveTo")
	}

	first := gd.NextPoint()
	cmd, count, err = gd.cmdAndCount()
	if err != nil {
		return nil, err
	}

	if cmd != lineTo {
		return nil, errors.New("second command not a lineTo")
	}

	ls := make(orb.LineString, 0, count+1)
	ls = append(ls, first)

	for i := uint32(0); i < count; i++ {
		ls = append(ls, gd.NextPoint())
	}

	return ls, nil
}

func (gd *geomDecoder) decodeLineString() (orb.Geometry, error) {
	var mls orb.MultiLineString
	for !gd.done() {
		ls, err := gd.decodeLine()
		if err != nil {
			return nil, err
		}

		if gd.done() && len(mls) == 0 {
			return ls, nil
		}

		mls = append(mls, ls)
	}

	return mls, nil
}

func (gd *geomDecoder) decodePolygon() (orb.Geometry, error) {
	var mp orb.MultiPolygon
	var p orb.Polygon
	for !gd.done() {
		ls, err := gd.decodeLine()
		if err != nil {
			return nil, err
		}

		r := orb.Ring(ls)

		cmd, _, err := gd.cmdAndCount()
		if err != nil {
			return nil, err
		}

		if cmd == closePath && !r.Closed() {
			r = append(r, r[0])
		}

		// figure out if new polygon
		if len(mp) == 0 && len(p) == 0 {
			p = append(p, r)
		} else {
			if r.Orientation() == orb.CCW {
				mp = append(mp, p)
				p = orb.Polygon{r}
			} else {
				p = append(p, r)
			}
		}
	}

	if len(mp) == 0 {
		return p, nil
	}

	return append(mp, p), nil
}

func (gd *geomDecoder) cmdAndCount() (uint32, uint32, error) {
	if gd.i >= len(gd.geom) {
		return 0, 0, errors.New("no more data")
	}

	v := gd.geom[gd.i]

	cmd := v & 0x07
	count := v >> 3
	gd.i++

	if cmd != closePath {
		if v := gd.i + int(2*count); len(gd.geom) < v {
			return 0, 0, errors.Errorf("data cut short: needed %d, have %d", v, len(gd.geom))
		}
	}

	return cmd, count, nil
}

func (gd *geomDecoder) NextPoint() orb.Point {
	gd.i += 2
	gd.prev[0] += unzigzag(gd.geom[gd.i-2])
	gd.prev[1] += unzigzag(gd.geom[gd.i-1])
	return gd.prev
}

func (gd *geomDecoder) done() bool {
	return gd.i >= len(gd.geom)
}

func newKeyValueEncoder() *keyValueEncoder {
	return &keyValueEncoder{
		keyMap:   make(map[string]uint32),
		valueMap: make(map[interface{}]uint32),
	}
}

func (kve *keyValueEncoder) Key(s string) uint32 {
	if i, ok := kve.keyMap[s]; ok {
		return i
	}

	i := uint32(len(kve.Keys))
	kve.Keys = append(kve.Keys, s)
	kve.keyMap[s] = i

	return i
}

func (kve *keyValueEncoder) Value(v interface{}) (uint32, error) {
	// If a type is not comparable we can't figure out uniqueness in the hash,
	// we also can't encode it into a vectortile.Tile_Value.
	// So we encoded it as a json string, which is what other encoders
	// also do.
	if !reflect.TypeOf(v).Comparable() {
		data, err := json.Marshal(v)
		if err != nil {
			return 0, errors.Errorf("uncomparable: %T", v)
		}

		v = string(data)
	}

	if i, ok := kve.valueMap[v]; ok {
		return i, nil
	}

	tv, err := encodeValue(v)
	if err != nil {
		return 0, err
	}

	i := uint32(len(kve.Values))
	kve.Values = append(kve.Values, tv)
	kve.valueMap[v] = i

	return i, nil
}

func encodeValue(v interface{}) (*vectortile.Tile_Value, error) {
	tv := &vectortile.Tile_Value{}
	switch t := v.(type) {
	case string:
		tv.StringValue = &t
	case fmt.Stringer:
		s := t.String()
		tv.StringValue = &s
	case int:
		i := int64(t)
		tv.SintValue = &i
	case int8:
		i := int64(t)
		tv.SintValue = &i
	case int16:
		i := int64(t)
		tv.SintValue = &i
	case int32:
		i := int64(t)
		tv.SintValue = &i
	case int64:
		i := int64(t)
		tv.SintValue = &i
	case uint:
		i := uint64(t)
		tv.UintValue = &i
	case uint8:
		i := uint64(t)
		tv.UintValue = &i
	case uint16:
		i := uint64(t)
		tv.UintValue = &i
	case uint32:
		i := uint64(t)
		tv.UintValue = &i
	case uint64:
		i := uint64(t)
		tv.UintValue = &i
	case float32:
		tv.FloatValue = &t
	case float64:
		tv.DoubleValue = &t
	case bool:
		tv.BoolValue = &t
	default:
		return nil, errors.Errorf("unable to encode value of type %T: %v", v, v)
	}

	return tv, nil
}

func decodeValue(v *vectortile.Tile_Value) interface{} {
	if v == nil {
		return nil
	}

	if v.StringValue != nil {
		return *v.StringValue
	} else if v.FloatValue != nil {
		return float64(*v.FloatValue)
	} else if v.DoubleValue != nil {
		return *v.DoubleValue
	} else if v.IntValue != nil {
		return float64(*v.IntValue)
	} else if v.UintValue != nil {
		return float64(*v.UintValue)
	} else if v.SintValue != nil {
		return float64(*v.SintValue)
	} else if v.BoolValue != nil {
		return *v.BoolValue
	}

	return nil
}

// functions to estimate encoded length

func elMLS(mls orb.MultiLineString) int {
	c := 0
	for _, ls := range mls {
		c += 2 + 2*len(ls)
	}

	return c
}

func elP(p orb.Polygon) int {
	c := 0
	for _, r := range p {
		c += 3 + 2*len(r)
	}

	return c
}

func elMP(mp orb.MultiPolygon) int {
	c := 0
	for _, p := range mp {
		c += elP(p)
	}

	return c
}

func unzigzag(v uint32) float64 {
	return float64(int32(((v >> 1) & ((1 << 32) - 1)) ^ -(v & 1)))
}
