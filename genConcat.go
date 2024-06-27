package revregex

import (
	"context"
	"regexp/syntax"
)

type genConcat2 struct {
	// immutable
	ctx  context.Context
	sub1 *syntax.Regexp // the first expression, garanteed non nil
	sub2 *syntax.Regexp // the second expression, garanteed non nil

	// state management, use Reset to set initially
	len1, len2 int       // both lengths, initially n,0
	gen1, gen2 Generator // the generators, initially nil before reset
	frag1      *string   // fragment already read from gen1, nil if not read yet
}

var _ Generator = new(genConcat2)

// Reset implements Generator.
func (g *genConcat2) Reset(exactLength int) (err error) {
	if exactLength < 0 {
		return ErrInvalidLength
	}

	// reset the lengths and the generators
	g.len1, g.len2 = exactLength+1, -1
	err = g.incSplit()
	if err != nil {
		return err
	}

	// reset frag1
	g.frag1 = nil

	// all good
	return g.ctx.Err()
}

// Next implements Generator.
func (g *genConcat2) Next() (f string, err error) {

	// set frag1 if not set yet, set it
	if g.frag1 == nil {
		g.frag1 = new(string)
		*g.frag1, err = g.gen1.Next()
		if err != nil {
			// try to change split
			err = g.incSplit()
			if err != nil {
				return "", ErrDone
			}
			// iterate
			return g.Next()
		}
	}

	// generate frag2
	f2, err := g.gen2.Next()
	if err != nil {
		// try to change split
		err = g.incSplit()
		if err != nil {
			return "", ErrDone
		}
		// iterate
		return g.Next()
	}

	// return concatenated frags
	return *g.frag1 + f2, g.ctx.Err()

}

// increment split. Return error if no more split to try.
func (g *genConcat2) incSplit() (err error) {
	if g.len1 == 0 {
		return ErrDone
	}
	// try with a new split
	g.len1 = g.len1 - 1
	g.len2 = g.len2 + 1

	// reset frag1
	g.frag1 = nil

	// set the Generators
	g.gen1, err = newGenerator(g.ctx, g.sub1, g.len1)
	if err != nil {
		return err
	}
	g.gen2, err = newGenerator(g.ctx, g.sub2, g.len2)
	if err != nil {
		return err
	}

	// return success
	return g.ctx.Err()
}
