package rgen

import (
	"context"
	"fmt"
	"regexp/syntax"
)

var ErrDone = fmt.Errorf("done")

type Generator interface {

	// error only on irrecoverable erorrs.
	// never error if nothing can be generated.
	Reset(length int) error

	// Generate nex string
	// returns error if no more string available.
	Next() error

	// retrieve the last generated string.
	// Undefined if called before Next() is called.
	Last() string
}

// Compile a pattern into a new generator.
// You must provide the maximum length for further resets.
// The new generator is already reset.
func NewGenerator(ctx context.Context, pattern string, max int) (Generator, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if max < 0 {
		return nil, fmt.Errorf("max length must be >= 0")
	}
	if pattern == "" {
		return nil, fmt.Errorf("pattern must not be empty")
	}
	re, err := syntax.Parse(pattern, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	return newGenerator(ctx, re, max)
}

// generic constructor from sub regex tree
func newGenerator(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {

	re = re.Simplify()

	switch re.Op {
	case syntax.OpNoMatch:
		return &genNoMatch{}, ctx.Err()
	case syntax.OpEmptyMatch:
		return &genEmptyMatch{}, ctx.Err()
	case syntax.OpLiteral:
		return newGenLiteral(ctx, re, max)
	case syntax.OpAlternate:
		return newGenAlternate(ctx, re, max)
	case syntax.OpCapture:
		return newGenerator(ctx, re.Sub0[0], max)
	case syntax.OpCharClass:
		return newGenClass(ctx, re, max)
	case syntax.OpQuest:
		return newGenQuest(ctx, re, max)
	case syntax.OpConcat:
		return newGenConcat(ctx, re, max)
	case syntax.OpStar:
		return newGenStar(ctx, re, max)
	default:
		panic(fmt.Sprintf("unknown operation : %d", re.Op))
	}
}

// ====================================================

// common structure for basic generator implementation
type generator struct {
	ctx  context.Context
	max  int
	gens []Generator
	last string
	done bool
}

var _ Generator = new(generator)

// Next implements Generator.
func (g *generator) Next() error {
	panic("unimplemented")
}

func (g *generator) Last() string {
	return g.last
}

// Reset implements Generator.
func (g *generator) Reset(length int) error {

	if length < 0 {
		return fmt.Errorf("length must be >= 0")
	}
	if length > g.max {
		return fmt.Errorf("length %d must be <= %d", length, g.max)
	}
	g.done = false
	g.last = ""
	for _, gen := range g.gens {
		if g.ctx.Err() != nil {
			return g.ctx.Err()
		}
		if err := gen.Reset(length); err != nil {
			return err
		}
	}
	return g.ctx.Err()
}

// ======================================================================

// match nothing
type genNoMatch struct{}

var _ Generator = new(genNoMatch)

// Next implements Generator.
func (g *genNoMatch) Next() error {
	return ErrDone
}

// Reset implements Generator.
func (g *genNoMatch) Reset(length int) error {
	return nil
}

func (g *genNoMatch) Last() string {
	panic("last should never be called on genNoMatch")
}

// ======================================================================

// only match ""
type genEmptyMatch struct {
	done bool
}

var _ Generator = new(genEmptyMatch)

// Next implements Generator.
func (g *genEmptyMatch) Next() error {
	if g.done {
		return ErrDone
	}
	g.done = true
	return nil
}

// Reset implements Generator.
func (g *genEmptyMatch) Reset(length int) error {
	g.done = (length != 0)
	return nil
}

func (g *genEmptyMatch) Last() string {
	return ""
}

// ======================================================================

// literal string
