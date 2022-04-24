package revregex

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"
)

// Gen can generate deterministic or random strings that match a given regexp.
// Gen is thread safe.
type Gen struct {
	// source string for the regexp.
	source string
	// root parsed tree
	tree *syntax.Regexp
}

// Same as NewGen, but in addition, the tree is simplified.
func NewGenSimpl(source string) *Gen {
	g := NewGen(source)
	g.tree = g.tree.Simplify()
	return g
}

// NewGen creates a new generator.
// It will panic if the regexp provided is not syntacly correct.
// Use POSIX syntax. No tree simplification.
func NewGen(source string) *Gen {
	var err error
	g := new(Gen)
	g.source = source
	g.tree, err = syntax.Parse(source, syntax.POSIX)
	if err != nil {
		panic(err)
	}
	return g
}

func (g *Gen) String() string {

	var b strings.Builder

	fmt.Fprintf(&b, "%q\t->\t", g.source)

	toString(&b, g.tree)
	return b.String()

}

func toString(b *strings.Builder, re *syntax.Regexp) {
	if re == nil {
		fmt.Fprint(b, nil)
		return
	}
	switch re.Op {

	case
		syntax.OpLiteral, syntax.OpCharClass:
		fmt.Fprintf(b, "%s(%q)", re.Op, re.Rune)

	case syntax.OpRepeat:
		fmt.Fprintf(b, "%s{%d,%d}(", re.Op, re.Min, re.Max)
		for _, rs := range re.Sub {
			toString(b, rs)
		}
		fmt.Fprint(b, ")")

	default:

		fmt.Fprint(b, re.Op, "(")
		for _, rs := range re.Sub {
			toString(b, rs)
		}
		fmt.Fprint(b, ")")
	}
}

var ErrVerificationFailed = fmt.Errorf("verification failed")

// Verify if a string match the regexp used to create g.
func (g *Gen) Verify(s string) error {

	ok, err := regexp.Match(g.source, []byte(s))

	if err != nil {
		return err
	}
	if !ok {
		return ErrVerificationFailed
	}
	return nil
}

// Next generate a new string that match the provided regexp.
func (g *Gen) Next(it Inter) string {
	var b strings.Builder
	next(&b, it, g.tree)
	return b.String()
}

func next(b *strings.Builder, it Inter, re *syntax.Regexp) {

	if re == nil {
		return
	}

	switch re.Op {
	case syntax.OpNoMatch: // matches no strings
		panic(re.Op.String() + " is not implemented")
	case syntax.OpEmptyMatch: // matches empty string
		return
	case syntax.OpLiteral: // matches Runes sequence
		fmt.Fprintf(b, "%s", string(re.Rune))
		return
	case syntax.OpCharClass: // matches Runes interpreted as range pair list
		// count choices ?
		nn := 0
		for i := 0; i+1 < len(re.Rune); i += 2 {
			nn += int(re.Rune[i+1]-re.Rune[i]) + 1
		}
		if nn == 0 {
			return
		}
		n := it.Intn(nn)
		for i := 0; i+1 < len(re.Rune); i += 2 {
			if n < int(re.Rune[i+1]-re.Rune[i])+1 { // match this pair !
				fmt.Fprintf(b, "%c", re.Rune[i]+rune(n))
				return // done !
			} else {
				// adjust n and continue to next pair
				n = n - (int(re.Rune[i+1]-re.Rune[i]) + 1)
			}
		}
		panic("internal error OpCharClass")
	case syntax.OpAnyCharNotNL: // matches any character except newline
		n := uint(it.Intn('\U0010ffff' - 1))
		if n == '\n' {
			n++
		}
		fmt.Fprintf(b, "%c", rune(n))
		return
	case syntax.OpAnyChar: // matches any character
		n := uint(it.Intn('\U0010ffff'))
		fmt.Fprintf(b, "%c", rune(n))
		return
	case syntax.OpBeginLine: // matches empty string at beginning of line
		panic(re.Op.String() + " is not implemented")
	case syntax.OpEndLine: // matches empty string at end of line
		panic(re.Op.String() + " is not implemented")
	case syntax.OpBeginText: // matches empty string at beginning of text
		panic(re.Op.String() + " is not implemented")
	case syntax.OpEndText: // matches empty string at end of text
		panic(re.Op.String() + " is not implemented")
	case syntax.OpWordBoundary: // matches word boundary `\b`
		panic(re.Op.String() + " is not implemented")
	case syntax.OpNoWordBoundary: // matches word non-boundary `\B`
		panic(re.Op.String() + " is not implemented")
	case syntax.OpCapture: // capturing subexpression with index Cap, optional name Name
		for _, sub := range re.Sub {
			next(b, it, sub)
		}
		return
	case syntax.OpStar: // matches Sub[0] zero or more times
		n := exp(it)
		for i := 0; i < n; i++ {
			for _, sub := range re.Sub { // the choice may differ between the repetitions !
				next(b, it, sub)
			}
		}
	case syntax.OpPlus: // matches Sub[0] one or more times
		n := exp(it) + 1
		for i := 0; i < n; i++ {
			for _, sub := range re.Sub { // the choice may differ between the repetitions !
				next(b, it, sub)
			}
		}
	case syntax.OpQuest: // matches Sub[0] zero or one times
		if it.Intn(2) == 0 {
			return
		} else {
			for _, sub := range re.Sub {
				next(b, it, sub)
			}
			return
		}
	case syntax.OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		n := re.Min + it.Intn(re.Max-re.Min+1)
		for i := 0; i < n; i++ {
			for _, sub := range re.Sub { // the choice may differ between the repetitions !
				next(b, it, sub)
			}
		}
	case syntax.OpConcat: // matches concatenation of Subs
		for _, sub := range re.Sub {
			next(b, it, sub)
		}
		return
	case syntax.OpAlternate: // matches alternation of Subs
		n := it.Intn(len(re.Sub))
		next(b, it, re.Sub[n])
		return
	default:
		panic("uniplemented regexp tree operation")
	}

}
