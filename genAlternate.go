package revregex

import (
	"context"
	"regexp/syntax"
)

type genAlternate struct {
	// immutable
	ctx  context.Context
	subs []*syntax.Regexp
	// state
	len int
	alt int       // alternate index to try next
	gen Generator // generator for current alternative
}

var _ Generator = new(genAlternate)

// Reset implements Generator.
func (g *genAlternate) Reset(exactLength int) error {
	if exactLength < 0 {
		return ErrInvalidLength
	}
	g.gen = nil
	g.alt = 0
	g.len = exactLength
	return g.ctx.Err()
}

// Next implements Generator.
func (g *genAlternate) Next() (f Fragment, err error) {

	if g.alt >= len(g.subs) {
		return Fragment{}, ErrDone
	}

	if g.gen == nil {
		g.gen, err = newGenerator(g.ctx, g.subs[g.alt], g.len)
		if err != nil {
			return Fragment{}, err
		}
	}

	f, err = g.gen.Next()
	if err == nil {
		return f, g.ctx.Err()
	}
	// try next alternative
	g.alt++
	g.gen = nil
	return g.Next()
}
