package wkb

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/paulmach/orb"
)

func readPolygon(r io.Reader, bom binary.ByteOrder) (orb.Polygon, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	if num > maxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		num = maxMultiAlloc
	}

	result := make(orb.Polygon, 0, num)
	for i := 0; i < int(num); i++ {
		ls, err := readLineString(r, bom)
		if err != nil {
			return nil, err
		}

		result = append(result, orb.Ring(ls))
	}

	return result, nil
}

func (e *Encoder) writePolygon(p orb.Polygon) error {
	e.order.PutUint32(e.buf, polygonType)
	e.order.PutUint32(e.buf[4:], uint32(len(p)))
	_, err := e.w.Write(e.buf[:8])
	if err != nil {
		return err
	}
	for _, r := range p {
		e.order.PutUint32(e.buf, uint32(len(r)))
		_, err := e.w.Write(e.buf[:4])
		if err != nil {
			return err
		}
		for _, p := range r {
			e.order.PutUint64(e.buf, math.Float64bits(p[0]))
			e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
			_, err = e.w.Write(e.buf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readMultiPolygon(r io.Reader, bom binary.ByteOrder) (orb.MultiPolygon, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	if num > maxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		num = maxMultiAlloc
	}

	result := make(orb.MultiPolygon, 0, num)
	for i := 0; i < int(num); i++ {
		byteOrder, typ, err := readByteOrderType(r)
		if err != nil {
			return nil, err
		}

		if typ != polygonType {
			return nil, errors.New("expect multipolygons to contains polygons, did not find a polygon")
		}

		p, err := readPolygon(r, byteOrder)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

func (e *Encoder) writeMultiPolygon(mp orb.MultiPolygon) error {
	e.order.PutUint32(e.buf, multiPolygonType)
	e.order.PutUint32(e.buf[4:], uint32(len(mp)))
	_, err := e.w.Write(e.buf[:8])
	if err != nil {
		return err
	}

	for _, p := range mp {
		err := e.Encode(p)
		if err != nil {
			return err
		}
	}

	return nil
}
