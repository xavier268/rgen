package rgen

import (
	"context"
	"regexp/syntax"
)

type genAlternate struct {
	*generator
	alt int // pointer to the alternative being generated
}

var _ Generator = new(genAlternate)

func newGenAlternate(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {
	g := &genAlternate{
		generator: &generator{
			ctx:  ctx,
			max:  max,
			gens: make([]Generator, len(re.Sub)),
		},
		alt: 0,
	}

	// create the sub generators
	var err error
	for i, sub := range re.Sub {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if g.gens[i], err = newGenerator(g.ctx, sub, g.max); err != nil {
			return nil, err
		}
	}

	// all good
	return g, ctx.Err()
}

func (g *genAlternate) Reset(n int) error {
	g.alt = 0
	return g.generator.Reset(n)
}

func (g *genAlternate) Next() (string, error) {
	if g.alt >= len(g.gens) {
		return "", ErrDone
	}

	s, err := g.gens[g.alt].Next()
	if err == nil {
		return s, g.ctx.Err()
	}

	// try next alternative
	g.alt++
	return g.Next()
}
