package rgen

import (
	"context"
	"regexp/syntax"
	"strings"
)

type genConcat struct {
	*generator
	lens    []int // length expected from each of the generators
	target  int   // target length to generate
	splitok bool  // has the current split first values been initialized already ?
}

func newGenConcat(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {

	if re.Op != syntax.OpConcat {
		panic("calling newGenConcat on an exprerssion that is not a concat op")
	}
	g := &genConcat{
		generator: &generator{
			ctx:  ctx,
			max:  max,
			gens: make([]Generator, len(re.Sub)),
			last: "",
			done: false,
		},
		lens:    make([]int, len(re.Sub)),
		target:  0,
		splitok: false,
	}
	for i, sub := range re.Sub {
		var err error
		g.gens[i], err = newGenerator(ctx, sub, max)
		if err != nil {
			return nil, err
		}
	}
	return g, ctx.Err()
}

// just the bare minimum. No string values are initialized.
func (g *genConcat) Reset(n int) error {
	if n < 0 {
		panic("illegal negative length at reset")
	}
	// reset lens and generators 0, 0, 0, 0, ...,n
	for i, gen := range g.gens {
		if i == len(g.lens)-1 {
			g.lens[i] = n
		} else {
			g.lens[i] = 0
		}
		if err := gen.Reset(g.lens[i]); err != nil {
			return err
		}
	}
	g.splitok = false
	g.target = n
	g.done = false
	g.last = ""
	return g.ctx.Err()
}

func (g *genConcat) Next() error {

	// fmt.Printf("DEBUG : calling NextConcat with %#v\n", g)
	// fmt.Printf("DEBUG : generator content:  %#v\n", g.generator)

	if g.done {
		return nil
	}

	if !g.splitok {
		err := g.useNewSplit()
		if err != nil {
			return err
		}
		g.splitok = true
		g.setLast()
		return g.ctx.Err()
	}

	// now, g.splitok = true
	// here we already have an initialized split, with valid values.
	// lets try to increment ?
	err := g.doNext(0)
	if err != nil {
		// we could not increment, we need another split.
		err := g.useNewSplit()
		if err != nil {
			return err
		}
		// ok - we have a new useable split.
		g.setLast()
		return g.ctx.Err()
	}
	// here, we could increment within the existing split.
	g.setLast()
	return g.ctx.Err()
}

// ====================utilities====================

// try to execute next, or return error
// this will never change the split.
// It will only look at gens[i:].
// g.last is not updated.
// it assumes all first values of the split have been initialized.
func (g *genConcat) doNext(i int) error {

	if i >= len(g.gens) {
		return ErrDone
	}

	// try to use i
	if g.gens[i].Next() == nil {
		// success, return
		return nil
	} else {
		// reset i
		err := g.gens[i].Reset(g.lens[i])
		if err != nil {
			return err
		}
		// retrieve first value
		err = g.gens[i].Next()
		if err != nil {
			return err
		}
	}

	if g.ctx.Err() != nil {
		return g.ctx.Err()
	}

	// since we could not use i, try to increment i+1
	return g.doNext(i + 1)
}

// loop until a new valild split can be found.
// reset is performed, and first values are initilized.
// error means no further split will ever work.
// if slitok is false, the current split is just reset and initialized.
func (g *genConcat) useNewSplit() error {
	if g.splitok { // if current split was initilized, look for another one.
		err := incSplitConcat(g.lens)
		if err != nil {
			// no more split to try
			return ErrDone
		}
	}
	// here, we do not increment the split, but first try to use it.
	for i, gen := range g.gens {

		if g.ctx.Err() != nil {
			return g.ctx.Err()
		}

		if err := gen.Reset(g.lens[i]); err != nil {
			g.splitok = true       // force incrementing split
			return g.useNewSplit() // cannot reset this generator - try another split
		}
		if err := gen.Next(); err != nil {
			g.splitok = true       // force incrementing split
			return g.useNewSplit() // cannot initialize first value try another split
		}
	}
	// all good
	g.splitok = true
	return g.ctx.Err()

}

// sum a slice of ints
func sum(ll []int) int {
	sum := 0
	for _, v := range ll {
		sum += v
	}
	return sum
}

// increment a split for concat, starting from 0, 0, 0, ...., n
// Sum of values should always stay the same.
// all values are positive or 0.
func incSplitConcat(ll []int) error {
	if len(ll) <= 1 {
		return ErrDone // nothing to split
	}
	lasti := len(ll) - 1 // last index

	for i := 0; i < lasti; i++ {
		if ll[lasti] > 0 {
			ll[i]++
			ll[lasti]--
			return nil
		} else {
			// prepare to try next index ...
			ll[lasti] = ll[lasti] + ll[i]
			ll[i] = 0
			continue
		}
	}

	return ErrDone
}

// collect all the 'last' values from the generators into a single last value.
func (g *genConcat) setLast() {
	buf := new(strings.Builder)
	for _, gen := range g.gens {
		buf.WriteString(gen.Last())
	}
	g.last = buf.String()
}
