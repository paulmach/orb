package clip

import (
	"errors"

	"github.com/paulmach/orb"
)

// AroundBound will connect the endpoints of the linestring provided
// by wrapping the line around the bounds in the direction provided.
// Will append to the input.
func AroundBound(
	box Bound,
	in LineString,
	o orb.Orientation,
	orientation func(LineString) orb.Orientation,
) (LineString, error) {
	next, ok := nexts[o]
	if !ok {
		return nil, errors.New("wrap: invalid orientation")
	}

	fX, fY := in.Get(0)
	lX, lY := in.Get(in.Len() - 1)

	if fX == lX && fY == lY {
		return in, nil // endpoints match
	}

	target := box.bitCodeOpen(fX, fY)
	current := box.bitCodeOpen(lX, lY)

	if target == 0 || current == 0 {
		return in, errors.New("wrap: endpoints must be outside bound")
	}

	if current == target && orientation(in) == o {
		in.Append(fX, fY)
		return in, nil
	}

	// move to next and go until we're all the way around.
	current = next[current]
	for target != current {
		in.Append(box.pointFor(current))
		current = next[current]
	}

	// add first point to the end to make it a ring
	in.Append(fX, fY)
	return in, nil
}

//         left  mid  right
//    top  1001  1000  1010
//    mid  0001  0000  0010
// bottom  0101  0100  0110

// on the boundary is outside
func (b Bound) bitCodeOpen(x, y float64) int {
	code := 0
	if x <= b.Left {
		code |= 1
	} else if x >= b.Right {
		code |= 2
	}

	if y <= b.Bottom {
		code |= 4
	} else if y >= b.Top {
		code |= 8
	}

	return code
}

//         left  mid  right
//    top     9     8    10
//    mid     1     0     2
// bottom     5     4     6

// nexts takes a bitcode index and jumps to the next corner.
var nexts = map[orb.Orientation][11]int{
	orb.CW: [11]int{
		-1,
		9, // 1
		6, // 2
		-1,
		5, // 4
		1, // 5
		4, // 6
		-1,
		10, // 8
		8,  // 9
		2,  // 10
	},
	orb.CCW: [11]int{
		-1,
		5,  // 1
		10, // 2
		-1,
		6, // 4
		4, // 5
		2, // 6
		-1,
		9, // 8
		1, // 9
		8, // 10
	},
}
