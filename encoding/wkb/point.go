package wkb

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/paulmach/orb"
)

func readPoint(r io.Reader, order byteOrder, buf []byte) (orb.Point, error) {
	var p orb.Point

	for i := 0; i < 2; i++ {
		if _, err := io.ReadFull(r, buf); err != nil {
			return orb.Point{}, err
		}
		if order == littleEndian {
			p[i] = math.Float64frombits(binary.LittleEndian.Uint64(buf))
		} else {
			p[i] = math.Float64frombits(binary.BigEndian.Uint64(buf))
		}
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

func readMultiPoint(r io.Reader, order byteOrder, buf []byte) (orb.MultiPoint, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > maxPointsAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = maxPointsAlloc
	}
	result := make(orb.MultiPoint, 0, alloc)

	for i := 0; i < int(num); i++ {
		pOrder, typ, err := readByteOrderType(r, buf)
		if err != nil {
			return nil, err
		}

		if typ != pointType {
			return nil, errors.New("expect multipoint to contains points, did not find a point")
		}

		p, err := readPoint(r, pOrder, buf)
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
