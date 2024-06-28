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
	Next() (f string, err error)
}

// Generator will produce all matching strings for given pattern and length.
// It will return error if pattern or length is invalid.
// Even if generator will match nothing, there should not be an error.
// Calling Next() will only return strings with the EXACT requiered length.
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
		g := &genCClass{
			ctx:  ctx,
			min:  []rune{},
			max:  []rune{},
			r:    0,
			i:    0,
			done: false,
		}
		if len(re.Rune) == 0 {
			g.min = []rune{re.Rune0[0]}
			g.max = []rune{re.Rune0[1]}
		}
		for i := 0; i < len(re.Rune); i += 2 {
			g.min = append(g.min, re.Rune[i])
			g.max = append(g.max, re.Rune[i+1])
		}
		err := g.Reset(length)
		if err != nil {
			return nil, err
		}
		return g, g.ctx.Err()

	default:
		return nil, fmt.Errorf("unsupported op: %v", re.Op)
	}

}

// Function generate will generate strings matching the pattern, and send them to the provided channel.
// The function will block until either all strings have been sent, or the context is cancelled.
// It is the caller responsability to close the channel when the function returns.
// All strings sent to channel will have the exact required length.
func Generate(ctx context.Context, pattern string, length int, out chan<- string) error {
	if out == nil {
		return fmt.Errorf("nil channel")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	g, err := NewGenerator(ctx, pattern, length)
	if err != nil {
		return err
	}
	s := ""
	for err == nil {
		s, err = g.Next()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err == nil {
				out <- s
			}
		}
		// loop until error becomes non nil
	}

	return ctx.Err()
}
