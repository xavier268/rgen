package revregex

import (
	"context"
)

type genLiteral struct {
	ctx  context.Context
	s    string
	done bool // do we have something to send ?
}

var _ Generator = new(genLiteral)

// Next implements Generator.
func (g *genLiteral) Next() (f string, err error) {
	if g.done {
		return "", ErrDone
	} else {
		g.done = true
		return g.s, g.ctx.Err()
	}
}

// Reset implements Generator.
func (g *genLiteral) Reset(n int) error {

	if n < 0 {
		return ErrInvalidLength
	}
	if n != len(g.s) {
		g.done = true
	} else {
		g.done = false
	}
	return g.ctx.Err()
}
