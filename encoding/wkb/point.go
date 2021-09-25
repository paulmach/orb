package wkb

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/paulmach/orb"
)

func unmarshalPoints(order binary.ByteOrder, data []byte) ([]orb.Point, error) {
	if len(data) < 4 {
		return nil, ErrNotWKB
	}
	num := order.Uint32(data)
	data = data[4:]

	if len(data) < int(num*16) {
		return nil, ErrNotWKB
	}

	alloc := num
	if alloc > maxPointsAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = maxPointsAlloc
	}
	result := make([]orb.Point, 0, alloc)

	for i := 0; i < int(num); i++ {
		result = append(result, orb.Point{
			math.Float64frombits(order.Uint64(data[16*i:])),
			math.Float64frombits(order.Uint64(data[16*i+8:])),
		})
	}

	return result, nil
}

func unmarshalPoint(order binary.ByteOrder, buf []byte) (orb.Point, error) {
	if len(buf) < 16 {
		return orb.Point{}, ErrNotWKB
	}

	var p = orb.Point{
		math.Float64frombits(order.Uint64(buf)),
		math.Float64frombits(order.Uint64(buf[8:])),
	}

	return p, nil
}

func readPoint(r io.Reader, order binary.ByteOrder, buf []byte) (orb.Point, error) {
	var p orb.Point

	for i := 0; i < 2; i++ {
		if _, err := io.ReadFull(r, buf); err != nil {
			return orb.Point{}, err
		}
		p[i] = math.Float64frombits(order.Uint64(buf))
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

func unmarshalMultiPoint(order binary.ByteOrder, data []byte) (orb.MultiPoint, error) {
	if len(data) < 4 {
		return nil, ErrNotWKB
	}
	num := order.Uint32(data)
	data = data[4:]

	alloc := num
	if alloc > maxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = maxMultiAlloc
	}
	result := make(orb.MultiPoint, 0, alloc)

	for i := 0; i < int(num); i++ {
		p, err := scanPoint(data)
		if err != nil {
			return nil, err
		}

		data = data[21:]
		result = append(result, p)
	}

	return result, nil
}

func readMultiPoint(r io.Reader, order binary.ByteOrder, buf []byte) (orb.MultiPoint, error) {
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
