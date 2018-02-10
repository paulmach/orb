package wkb

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/paulmach/orb"
)

func readPoint(r io.Reader, bom binary.ByteOrder) (orb.Point, error) {
	var p orb.Point

	if err := binary.Read(r, bom, &p[0]); err != nil {
		return orb.Point{}, err
	}

	if err := binary.Read(r, bom, &p[1]); err != nil {
		return orb.Point{}, err
	}

	return p, nil
}

func (e *Encoder) writePoint(p orb.Point) error {
	e.order.PutUint32(e.buf, pointType)
	_, err := e.w.Write(e.buf[:4])
	if err != nil {
		return err
	}

	e.order.PutUint64(e.buf, math.Float64bits(p[0]))
	e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
	_, err = e.w.Write(e.buf)
	return err
}

func readMultiPoint(r io.Reader, bom binary.ByteOrder) (orb.MultiPoint, error) {
	var num uint32 // Number of points.
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	if num > maxPointsAlloc {
		// invalid data can come in here and allocate tons of memory.
		num = maxPointsAlloc
	}

	result := make(orb.MultiPoint, 0, num)
	for i := 0; i < int(num); i++ {
		byteOrder, typ, err := readByteOrderType(r)
		if err != nil {
			return nil, err
		}

		if typ != pointType {
			return nil, errors.New("expect multipoint to contains points, did not find a point")
		}

		p, err := readPoint(r, byteOrder)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

func (e *Encoder) writeMultiPoint(mp orb.MultiPoint) error {
	e.order.PutUint32(e.buf, multiPointType)
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
