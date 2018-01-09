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

	result := make(orb.Collection, 0, num)
	for i := 0; i < int(num); i++ {
		geom, err := Read(r)
		if err != nil {
			return nil, err
		}

		result = append(result, geom)
	}

	return result, nil
}
