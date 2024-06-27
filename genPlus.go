package revregex

import (
	"context"
	"regexp/syntax"
)

type genPlus struct {
	// immutable
	ctx context.Context
	re  *syntax.Regexp
	// state management - replacing "x+" by "xx*"
	gen Generator
}

var _ Generator = new(genPlus)

// Reset implements Generator.
func (g *genPlus) Reset(exactLength int) (err error) {
	nt := &syntax.Regexp{
		Op: syntax.OpConcat,
		Sub: []*syntax.Regexp{
			g.re,
			{Op: syntax.OpStar, Sub0: [1]*syntax.Regexp{g.re}},
		},
	}
	g.gen, err = newGenerator(g.ctx, nt, exactLength)
	if err != nil {
		return err
	}
	return g.ctx.Err()
}

// Next implements Generator.
func (g *genPlus) Next() (f Fragment, err error) {
	return g.gen.Next()
}
