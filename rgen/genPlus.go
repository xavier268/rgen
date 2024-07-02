package rgen

import (
	"context"
	"fmt"
	"regexp/syntax"
)

func newGenPlus(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {
	if re == nil {
		panic("unexpected nil Regexp")
	}
	if re.Op != syntax.OpPlus {
		panic(fmt.Sprintf("invalid opcode to newGenPlus : %d", re.Op))
	}

	nt := &syntax.Regexp{
		Op: syntax.OpConcat, // xx*
		Sub: []*syntax.Regexp{
			re.Sub0[0], // x
			{ // x*
				Op:   syntax.OpStar,
				Sub0: re.Sub0, // x*
			},
		},
	}

	return newGenConcat(ctx, nt, max)
}
