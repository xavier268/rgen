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
// max specifiy the maximum length available for this generator.
func newGenLiteral(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {
	g := &genLiteral{
		generator: &generator{
			ctx:  ctx,
			max:  max,
			gens: nil,
			done: false,
		},
		value: string(re.Rune),
	}
	return g, ctx.Err()
}

// Next implements Generator.
func (g *genLiteral) Next() (string, error) {
	if g.done {
		return "", ErrDone
	}
	g.done = true
	return g.value, g.ctx.Err()
}

// Reset implements Generator.
func (g *genLiteral) Reset(n int) error {
	g.generator.Reset(n)
	g.done = !(n == len(g.value))
	return g.ctx.Err()
}
