package revregex

import (
	"context"
	"regexp/syntax"
)

type genQuest struct {
	ctx       context.Context
	re        *syntax.Regexp
	len       int
	emptyDone bool // is the empty alternative already sent ?
	gen       Generator
}

var _ Generator = new(genQuest)

// Reset implements Generator.
func (g *genQuest) Reset(exactLength int) (err error) {
	g.len = exactLength
	g.emptyDone = false
	g.gen, err = newGenerator(g.ctx, g.re, g.len)
	if err != nil {
		return err
	}
	return g.ctx.Err()
}

// Next implements Generator.
func (g *genQuest) Next() (f Fragment, err error) {
	if !g.emptyDone && g.len == 0 {
		g.emptyDone = true
		return Fragment{}, g.ctx.Err()
	}
	return g.gen.Next()
}
