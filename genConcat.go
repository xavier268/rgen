package revregex

import (
	"context"
	"regexp/syntax"
)

type genConcat2 struct {
	// immutable
	ctx  context.Context
	sub1 *syntax.Regexp // the first expression, garanteed non nul
	sub2 *syntax.Regexp // the second expression, garanteed non nul
	// state management
	len1, len2 int       // both lengths, initially n,0
	gen1, gen2 Generator // the generators, initially nil before reset
	frag1      Fragment  // fragment already read from gen1
	done       bool      // we are done, no more fragments will be generated !
}

var _ Generator = new(genConcat2)

// Reset implements Generator.
func (g *genConcat2) Reset(exactLength int) (err error) {
	if exactLength < 0 {
		return ErrInvalidLength
	}
	// reset the lengths
	g.len1, g.len2 = exactLength, 0
	// Set the Generators
	g.gen1, err = newGenerator(g.ctx, g.sub1, g.len1)
	if err != nil {
		return err
	}
	g.gen2, err = newGenerator(g.ctx, g.sub2, g.len2)
	if err != nil {
		return err
	}
	g.done = false

	// Set frag1 to the first fragment of gen1
	g.frag1, err = g.gen1.Next()
	if err != nil {
		// TODO - WRONG !!
		g.done = true // nothing will be generated !
	}

	// all good
	return g.ctx.Err()
}

// Next implements Generator.
func (g *genConcat2) Next() (f Fragment, err error) {
	if g.done {
		return Fragment{}, ErrDone
	}

	frag2, err := g.gen2.Next()
	if err == nil {
		fr3, err := concatFrags(g.frag1, frag2)
		if err == nil {
			return fr3, nil
		} else {
			// concat error, try next value
			return g.Next()
		}
	} else {
		// no more gen2 frags, change frag1 and reset gen2
		g.frag1, err = g.gen1.Next()
		if err == nil {
			g.gen2.Reset(g.len2)
			return g.Next()
		} else {
			// no more gen1 frags, we need to change split
			if g.len1 == 0 {
				// we're done, no more split can be tried
				return Fragment{}, ErrDone
			}
			// try with a new split
			g.len1 = g.len1 - 1
			g.len2 = g.len2 + 1

		}
	}

	panic("TODO")
}

func concatFrags(f1, f2 Fragment) (f3 Fragment, err error) {
	f3 = Fragment{
		s:     f1.s + f2.s,
		start: f1.start,
		end:   f2.end,
	}
	if f1.end && len(f2.s) > 0 {
		return f3, ErrConcatFragments
	}
	if f2.start && len(f1.s) > 0 {
		return f3, ErrConcatFragments
	}
	return f3, nil
}
