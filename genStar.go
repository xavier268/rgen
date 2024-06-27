package revregex

import (
	"context"
	"regexp/syntax"
)

type genStar struct {
	// immutable
	ctx context.Context
	re  *syntax.Regexp
	// state management
	len       int
	emptyDone bool      // used if len is 0, to remeber if 'nothing' was already sent.
	gen       Generator // a single generator, constructed using the length n, using  x?x?x?... (n times)
}

var _ Generator = new(genStar)

func (g *genStar) Reset(len int) (err error) {
	g.len = len
	if g.len == 0 {
		// if len != 0, never consider the 'nothing' solution.
		g.emptyDone = false
		g.gen = nil // do not use generator in this case
		return g.ctx.Err()
	}
	// Here, len > 0
	g.emptyDone = true // if len >0, empty should never be considered for the full solution

	// reconstruct a new tree,  x -> x?x?x?x? ...
	subs := make([]*syntax.Regexp, g.len)
	for i := range subs {
		subs[i] = &syntax.Regexp{
			Op:   syntax.OpQuest,
			Sub0: [1]*syntax.Regexp{g.re},
		} // x?
	}
	nt := &syntax.Regexp{
		Op:  syntax.OpConcat,
		Sub: subs,
	}
	// make generator for x?x?x? ...
	g.gen, err = newGenerator(g.ctx, nt, g.len)
	if err != nil {
		return err
	}

	// all good
	return g.ctx.Err()

}

func (g *genStar) Next() (Fragment, error) {

	if g.len == 0 && !g.emptyDone {
		g.emptyDone = true
		return Fragment{}, g.ctx.Err()
	}

	if g.len == 0 && g.emptyDone {
		return Fragment{}, ErrDone
	}

	return g.gen.Next()
}
