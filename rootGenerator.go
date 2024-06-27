package revregex

import (
	"context"
	"fmt"
	"regexp/syntax"
)

var ErrDone = fmt.Errorf("done")
var ErrInvalidLength = fmt.Errorf("invalid length")
var ErrEmptyConcat = fmt.Errorf("empty opConcat arguments")
var ErrConcatFragments = fmt.Errorf("incompatible fragements, concatenation is forbidden") // when trying to cancatenate with start/end constraints

var DEBUG = false

type Generator interface {
	// reset the generator for given length
	Reset(exactLength int) error
	// retrieve next matching string
	Next() (f Fragment, err error)
}

// Fragment are immutable
type Fragment struct {
	s     string // the string this fragment is based on
	start bool   // no prior fragment allowed
	end   bool   // no further fragment allowed
}

func (f Fragment) String() (s string) {
	s = f.s
	if f.start {
		s = "^" + s
	}
	if f.end {
		s = s + "$"
	}
	return s
}

// Generator will produce all matching strings for given pattern and length.
// It will return error if pattern or length is invalid.
func NewGenerator(ctx context.Context, pattern string, length int) (Generator, error) {
	re, err := syntax.Parse(pattern, syntax.POSIX)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	re = re.Simplify()

	if DEBUG {
		fmt.Printf("parsed simplified regexp : %s, length: %d\n", re.String(), length)
	}

	return newGenerator(ctx, re, length)
}

// newGenerator will create a generator from a regexp already compiled.
// Generator is already initialized.
func newGenerator(ctx context.Context, re *syntax.Regexp, length int) (Generator, error) {

	if DEBUG {
		fmt.Printf("calling newGenerator: %s, length: %d\n", re.String(), length)
	}

	// check length
	if length < 0 {
		return nil, ErrInvalidLength
	}

	// check context
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// select generator based on top op value
	switch re.Op {

	case syntax.OpLiteral: // xyz
		g := &genLiteral{
			ctx: ctx,
			s:   string(re.Rune),
		}
		err := g.Reset(length)
		return g, err

	case syntax.OpAlternate: // a | b
		g := &genAlternate{
			ctx:  ctx,
			subs: re.Sub,
		}
		err := g.Reset(length)
		return g, err

	case syntax.OpCapture: // (ab)
		return newGenerator(ctx, re.Sub0[0], length)

	case syntax.OpConcat: // (ab)(cd)
		if len(re.Sub) == 0 {
			return nil, ErrEmptyConcat
		}
		if len(re.Sub) == 1 {
			return newGenerator(ctx, re.Sub[0], length)
		}

		// here, we split between the first exp and the rest of the concat
		rest := &syntax.Regexp{
			Op:  syntax.OpConcat,
			Sub: re.Sub[1:], // non empty, but could be only 1
		}
		g := &genConcat2{
			ctx:  ctx,
			sub1: re.Sub[0],
			sub2: rest,
		}
		err := g.Reset(length)
		if err != nil {
			return nil, err
		}
		return g, ctx.Err()

	case syntax.OpQuest: // (exp)?
		g := &genQuest{
			ctx:       ctx,
			re:        re.Sub0[0],
			len:       length,
			emptyDone: false,
			gen:       nil,
		}
		err := g.Reset(length)
		if err != nil {
			return nil, err
		}
		return g, g.ctx.Err()

	case syntax.OpStar:
		g := &genStar{
			ctx: ctx,
			re:  re.Sub0[0],
			len: length,
		}
		err := g.Reset(length)
		if err != nil {
			return nil, err
		}
		return g, g.ctx.Err()

	case syntax.OpPlus:
		g := &genPlus{
			ctx: ctx,
			re:  re.Sub0[0],
		}
		err := g.Reset(length)
		if err != nil {
			return nil, err
		}
		return g, g.ctx.Err()

	case syntax.OpCharClass:
		panic("not implemented")

	default:
		return nil, fmt.Errorf("unsupported op: %v", re.Op)
	}

}

func concatFrags(f1, f2 Fragment) (f3 Fragment, err error) {
	f3 = Fragment{
		s:     f1.s + f2.s,
		start: f1.start,
		end:   f2.end,
	}
	if f1.end && len(f2.s) > 0 {
		return f3, ErrConcatFragments
	}
	if f2.start && len(f1.s) > 0 {
		return f3, ErrConcatFragments
	}
	return f3, nil
}
