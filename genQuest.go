package rgen

import (
	"context"
	"regexp/syntax"
)

func newGenQuest(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {

	// modify re, replacing x? by <noMatch>|x
	nt := &syntax.Regexp{
		Op: syntax.OpAlternate,
		Sub: []*syntax.Regexp{
			{Op: syntax.OpEmptyMatch},
			re.Sub0[0],
		},
	}
	return newGenerator(ctx, nt, max)
}
