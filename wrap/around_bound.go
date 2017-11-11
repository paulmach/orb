package wrap

import (
	"errors"

	"github.com/paulmach/orb"
)

// aroundBound will connect the endpoints of the linestring provided
// by wrapping the line around the bounds in the direction provided.
// Will append to the input.
func aroundBound(
	box orb.Bound,
	in orb.Ring,
	o orb.Orientation,
) (orb.Ring, error) {
	if o != orb.CCW && o != orb.CW {
		panic("invalid orientation")
	}

	if len(in) == 0 {
		return nil, nil
	}

	next := nexts[o]

	f := in[0]
	l := in[len(in)-1]

	if f == l {
		return in, nil // endpoints match
	}

	target := bitCodeOpen(box, f)
	current := bitCodeOpen(box, l)

	if target == 0 || current == 0 {
		return in, errors.New("wrap: endpoints must be outside bound")
	}

	if current == target && in.Orientation() == o {
		in = append(in, f)
		return in, nil
	}

	// move to next and go until we're all the way around.
	current = next[current]
	for target != current {
		in = append(in, pointFor(box, current))
		current = next[current]
	}

	// add first point to the end to make it a ring
	in = append(in, f)
	return in, nil
}

//         left  mid  right
//    top  1001  1000  1010
//    mid  0001  0000  0010
// bottom  0101  0100  0110

// on the boundary is outside
func bitCodeOpen(b orb.Bound, p orb.Point) int {
	code := 0
	if p[0] <= b.Left() {
		code |= 1
	} else if p[0] >= b.Right() {
		code |= 2
	}

	if p[1] <= b.Bottom() {
		code |= 4
	} else if p[1] >= b.Top() {
		code |= 8
	}

	return code
}

// pointFor returns a representative point for the side of the given bitCode.
func pointFor(b orb.Bound, code int) orb.Point {
	switch code {
	case 1:
		return orb.Point{b.Left(), (b.Top() + b.Bottom()) / 2}
	case 2:
		return orb.Point{b.Right(), (b.Top() + b.Bottom()) / 2}
	case 4:
		return orb.Point{(b.Right() + b.Left()) / 2, b.Bottom()}
	case 5:
		return orb.Point{b.Left(), b.Bottom()}
	case 6:
		return orb.Point{b.Right(), b.Bottom()}
	case 8:
		return orb.Point{(b.Right() + b.Left()) / 2, b.Top()}
	case 9:
		return orb.Point{b.Left(), b.Top()}
	case 10:
		return orb.Point{b.Right(), b.Top()}
	}

	panic("invalid code")
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
