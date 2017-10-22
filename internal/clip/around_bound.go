package clip

import (
	"errors"

	"github.com/paulmach/orb"
)

// AroundBound TODO
// Will modify the input by appending points to the end.
func AroundBound(
	box Bound,
	in LineString,
	o orb.Orientation,
	orientation func(LineString) orb.Orientation,
) (LineString, error) {
	next, ok := nexts[o]
	if !ok {
		return nil, errors.New("clip: invalid orientation")
	}

	fX, fY := in.Get(0)
	lX, lY := in.Get(in.Len() - 1)

	if fX == lX && fY == lY {
		return in, nil // endpoints match
	}

	target := box.bitCode(fX, fY)
	current := box.bitCode(lX, lY)

	if target == 0 || current == 0 {
		return in, errors.New("clip: endpoints must be outside bound")
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
