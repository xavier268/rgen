package old

import (
	"context"
	"fmt"
	"regexp/syntax"
	"testing"
)

func TestIncSplitVisual(t *testing.T) {

	g := &genConcat2{
		ctx:   context.Background(),
		sub1:  &syntax.Regexp{Op: syntax.OpLiteral, Rune: []rune("ab")},
		sub2:  &syntax.Regexp{Op: syntax.OpLiteral, Rune: []rune("kjh")},
		len1:  5,
		len2:  -1,
		gen1:  nil,
		gen2:  nil,
		frag1: "",
		done:  false,
	}

	var err error
	for err == nil {
		err = g.incSplit()
		fmt.Println("split : ", g.len1, g.len2)
	}
}
