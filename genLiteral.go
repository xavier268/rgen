package rgen

import (
	"context"
	"regexp/syntax"
)

type genLiteral struct {
	*generator
	value string
}

var _ Generator = new(genLiteral)

// create a new genLiteral from a parse tree.
// max specify the maximum length available for this generator.
func newGenLiteral(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {
	g := &genLiteral{
		generator: &generator{
			ctx:  ctx,
			max:  max,
			gens: nil,
			last: "",
			done: false,
		},
		value: string(re.Rune),
	}
	return g, ctx.Err()
}

// Next implements Generator.
func (g *genLiteral) Next() error {
	if g.done {
		return ErrDone
	}
	g.done = true
	g.last = g.value
	return g.ctx.Err()
}

// Reset implements Generator.
func (g *genLiteral) Reset(n int) error {
	g.generator.Reset(n)
	g.done = !(n == len(g.value))
	return g.ctx.Err()
}
