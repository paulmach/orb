package wkb

import (
	"encoding/binary"
	"io"

	"github.com/paulmach/orb"
)

func readCollection(r io.Reader, bom binary.ByteOrder) (orb.Collection, error) {
	var num uint32
	if err := binary.Read(r, bom, &num); err != nil {
		return nil, err
	}

	if num > maxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		num = maxMultiAlloc
	}

	result := make(orb.Collection, 0, num)
	for i := 0; i < int(num); i++ {
		geom, err := NewDecoder(r).Decode()
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
