package rgen

import (
	"context"
	"regexp/syntax"
)

type genConcat struct {
	*generator
	lens   []int // length expected from each of the generators
	target int   // target length to generate
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
		lens:   make([]int, len(re.Sub)),
		target: 0,
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

func (g *genConcat) Reset(n int) error {
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
	g.target = n
	g.done = false
	g.last = ""
	return g.ctx.Err()
}

func (g *genConcat) Next() error {
	if g.done {
		return nil
	}
	// try to retrieve with existing split

	return g.ctx.Err()
}

// ========================================

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
