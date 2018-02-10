package wkb

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/paulmach/orb"
)

func readLineString(r io.Reader, bom binary.ByteOrder) (orb.LineString, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	if num > maxPointsAlloc {
		// invalid data can come in here and allocate tons of memory.
		num = maxPointsAlloc
	}

	result := make(orb.LineString, 0, num)
	for i := 0; i < int(num); i++ {
		p, err := readPoint(r, bom)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

func (e *Encoder) writeLineString(ls orb.LineString) error {
	e.order.PutUint32(e.buf, lineStringType)
	e.order.PutUint32(e.buf[4:], uint32(len(ls)))
	_, err := e.w.Write(e.buf[:8])
	if err != nil {
		return err
	}

	for _, p := range ls {
		e.order.PutUint64(e.buf, math.Float64bits(p[0]))
		e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
		_, err = e.w.Write(e.buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func readMultiLineString(r io.Reader, bom binary.ByteOrder) (orb.MultiLineString, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	if num > maxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		num = maxMultiAlloc
	}

	result := make(orb.MultiLineString, 0, num)
	for i := 0; i < int(num); i++ {
		byteOrder, typ, err := readByteOrderType(r)
		if err != nil {
			return nil, err
		}

		if typ != lineStringType {
			return nil, errors.New("expect multilines to contains lines, did not find a line")
		}

		ls, err := readLineString(r, byteOrder)
		if err != nil {
			return nil, err
		}

		result = append(result, ls)
	}

	return result, nil
}

func (e *Encoder) writeMultiLineString(mls orb.MultiLineString) error {
	e.order.PutUint32(e.buf, multiLineStringType)
	e.order.PutUint32(e.buf[4:], uint32(len(mls)))
	_, err := e.w.Write(e.buf[:8])
	if err != nil {
		return err
	}

	for _, ls := range mls {
		err := e.Encode(ls)
		if err != nil {
			return err
		}
	}

	return nil
}
