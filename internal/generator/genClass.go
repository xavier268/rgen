package generator

import (
	"context"
	"regexp/syntax"
)

type genClass struct {
	*generator
	values []rune // list of all acceptable values
	i      int    // current position in string
}

func newGenClass(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {
	g := &genClass{
		generator: &generator{
			ctx:  ctx,
			max:  max,
			gens: nil,
			done: false,
		},
		values: make([]rune, 0, 20),
		i:      0,
	}

	// initialize acceptable runes values
	for i := 0; i < len(re.Rune); i += 2 {
		start := re.Rune[i]
		end := re.Rune[i+1]
		for j := start; j <= end; j++ {
			g.values = append(g.values, j)
		}
	}
	if len(re.Rune) == 0 {
		start := re.Rune0[0]
		end := re.Rune0[1]
		for j := start; j <= end; j++ {
			g.values = append(g.values, j)
		}
	}

	// all good
	return g, ctx.Err()
}

func (g *genClass) Reset(n int) error {
	g.i = 0
	err := g.generator.Reset(n)
	if err != nil {
		return err
	}
	g.done = (n != 1)
	return g.ctx.Err()
}

func (g *genClass) Next() error {
	if g.done {
		return ErrDone
	}

	if g.i >= len(g.values) {
		g.done = true
		return ErrDone
	}

	g.last = string(g.values[g.i])
	g.i++
	return g.ctx.Err()
}
