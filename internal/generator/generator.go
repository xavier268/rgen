// generator provides low level structures and functions used to efficiently generate strings from a regex pattern
package generator

import (
	"context"
	"fmt"
	"regexp/syntax"
)

var ErrDone = fmt.Errorf("done")

type Generator interface {

	// Sould be called before next is called.
	// Error only on irrecoverable errors.
	// Never error if nothing can be generated.
	Reset(length int) error

	// Generate next string
	// Returns error if no more string available.
	Next() error

	// Retrieve the last generated string.
	// May be called multiple times, and will always
	// return the same result until Next() is called again.
	// Undefined if called before Next() is called.
	Last() string
}

// Compile a pattern into a new Generator.
// You must provide a maximum length (ie, capacity) for further resets.
// The new generator has to be reset to the target length(less or equal to max)
// before Next() can be called and Last() can be read.
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

	re = re.Simplify()

	return newGenerator(ctx, re, max)
}

// generic constructor from sub regex tree
func newGenerator(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {

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
	case syntax.OpPlus:
		return newGenPlus(ctx, re, max)
	default:
		panic(fmt.Sprintf("opcode operation %d is not supported", re.Op))
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
