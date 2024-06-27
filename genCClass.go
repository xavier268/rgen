package revregex

import "context"

type genCClass struct {
	// immutable
	ctx context.Context
	min []rune // start of ranges
	max []rune // end of ranges, SHOULD BE EXACT SAME LENGTH as g.min !
	// state management
	r    int   // which range we are on
	i    int32 // index in range, int32 so that it can be a rune.
	done bool
}

var _ Generator = new(genCClass)

func (g *genCClass) Reset(exactLength int) error {
	g.r = 0
	g.i = 0
	g.done = (exactLength != 1) || len(g.min) == 0
	return g.ctx.Err()
}

func (g *genCClass) Next() (f Fragment, err error) {
	if g.done {
		return Fragment{}, ErrDone
	}

	// normalize pointers, error if cannot be normalized
	if g.i > g.max[g.r]-g.min[g.r] {
		g.i = 0
		g.r = g.r + 1
	}
	if g.r >= len(g.min) {
		g.done = true
		return Fragment{}, ErrDone
	}

	// fetch value currently pointed to by g.i,
	f.s = string(g.min[g.r] + g.i)
	f.start, f.end = false, false

	// increment (but do not normalize yet)
	g.i = g.i + 1

	// return value
	return f, g.ctx.Err()
}
