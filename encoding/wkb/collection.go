package wkb

import (
	"io"

	"github.com/paulmach/orb"
)

func readCollection(r io.Reader, order byteOrder, buf []byte) (orb.Collection, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > maxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = maxMultiAlloc
	}
	result := make(orb.Collection, 0, alloc)

	d := NewDecoder(r)
	for i := 0; i < int(num); i++ {
		geom, err := d.Decode()
		if err != nil {
			return nil, err
		}

		result = append(result, geom)
	}

	return result, nil
}

func (e *Encoder) writeCollection(c orb.Collection) error {
	e.order.PutUint32(e.buf, geometryCollectionType)
	e.order.PutUint32(e.buf[4:], uint32(len(c)))
	_, err := e.w.Write(e.buf[:8])
	if err != nil {
		return err
	}

	for _, geom := range c {
		err := e.Encode(geom)
		if err != nil {
			return err
		}
	}

	return nil
}
