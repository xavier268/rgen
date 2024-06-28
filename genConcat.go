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
	frag1      string    // fragment already read from gen1, nil if not read yet
	done       bool      // no more solutions
}

var _ Generator = new(genConcat2)

// Reset implements Generator.
func (g *genConcat2) Reset(exactLength int) (err error) {
	if exactLength < 0 {
		return ErrInvalidLength
	}
	g.done = false

	// reset the lengths and the generators
	g.len1, g.len2 = exactLength+1, -1
	err = g.incSplit() // set to split = (len, 0), and set frag1
	if err != nil {
		// there will be no solution !
		g.done = true
		return g.ctx.Err()
	}

	// all good
	return g.ctx.Err()
}

// Next implements Generator.
func (g *genConcat2) Next() (f string, err error) {

	if g.done {
		return "", ErrDone
	}

	// Frag1 is already set.

	// generate frag2
	f2, err := g.gen2.Next()
	if err == nil {
		// return concatenated frags
		return g.frag1 + f2, g.ctx.Err()
	}

	// here, we could not generate frag2.
	// increment frag1
	g.frag1, err = g.gen1.Next()
	if err != nil {
		// No more frag1, no more solutions
		g.done = true
		return "", ErrDone
	}
	// reset frag2 and loop
	err = g.gen2.Reset(g.len2)
	if err != nil {
		// Internal frag2 error, should have been detected at incSplit
		panic("internal error - frag2 reset impossible")
	}

	// loop
	return g.Next()
}

// increment split. Recreate both generators. Return error if no more split to try.
func (g *genConcat2) incSplit() (err error) {

	// test if more split are possible ?
	if g.done || g.len1 <= 0 {
		g.done = true
		return ErrDone
	}
	// try with a new split
	g.len1 = g.len1 - 1
	g.len2 = g.len2 + 1

	// set the Generators
	g.gen1, err = newGenerator(g.ctx, g.sub1, g.len1)
	if err != nil {
		return err
	}
	g.gen2, err = newGenerator(g.ctx, g.sub2, g.len2)
	if err != nil {
		return err
	}

	// initialize frag1
	g.frag1, err = g.gen1.Next()
	if err != nil {
		// try next split !
		return g.incSplit()
	}

	// return success
	return g.ctx.Err()
}
