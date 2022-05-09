package revregex

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"
)

const VERSION = "0.3.0"

// MaxUnicode is the maximum Unicode character that can be generated.
const MaxUnicode = '\U0010ffff'

// Gen can generate deterministic or random strings that match a given regexp.
// Gen is thread safe.
type Gen struct {
	// source string for the regexp.
	source string
	// root parsed tree
	tree *syntax.Regexp
}

// Same as NewGen, but in addition, the tree is simplified.
func NewGenSimpl(source string) (*Gen, error) {
	g, err := NewGen(source)
	g.tree = g.tree.Simplify()
	return g, err
}

// NewGen creates a new generator.
// It returns an error if the regexp provided is not syntaxicaly correct.
// Use POSIX syntax. No implicit parse tree simplification.
func NewGen(source string) (*Gen, error) {
	var err error
	g := new(Gen)
	g.source = cleanPattern(source)
	_, err = regexp.Compile(g.source)
	if err != nil {
		return nil, err
	}
	g.tree, err = syntax.Parse(g.source, syntax.POSIX)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Must is a utilty to panic on error, when creating a Gen.
// Typicla use is :
// g := Must(NewGen(pattern))
func Must(g *Gen, e error) *Gen {
	if e != nil {
		panic(e)
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

// Next generate a string that match the provided regexp, using the provided Chooser to make its choices.
func (g *Gen) Next(it Chooser) string {
	var b strings.Builder
	next(&b, it, g.tree)
	return b.String()
}

func next(b *strings.Builder, it Chooser, re *syntax.Regexp) {

	if re == nil {
		return
	}

	switch re.Op {
	case syntax.OpNoMatch, // matches no strings
		syntax.OpBeginLine,      // matches empty string at beginning of line
		syntax.OpEndLine,        // matches empty string at end of line
		syntax.OpBeginText,      // matches empty string at beginning of text
		syntax.OpEndText,        // matches empty string at end of text
		syntax.OpWordBoundary,   // matches word boundary `\b`
		syntax.OpNoWordBoundary: // matches word non-boundary `\B`
		var bb strings.Builder
		toString(&bb, re)
		panic(re.Op.String() + " is not implemented in : " + bb.String())
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
		n := uint(it.Intn(MaxUnicode - 1))
		if n == '\n' {
			n++
		}
		fmt.Fprintf(b, "%c", rune(n))
		return
	case syntax.OpAnyChar: // matches any character
		n := uint(it.Intn(MaxUnicode))
		fmt.Fprintf(b, "%c", rune(n))
		return

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
		panic("unimplemented regexp parse tree operation")
	}

}

// cleanPattern from valid but meaningless directives.
func cleanPattern(pattern string) string {

	for {
		pat := pattern
		pat = strings.ReplaceAll(pat, "^", "")
		pat = strings.ReplaceAll(pat, "$", "")
		pat = strings.ReplaceAll(pat, "\\A", "")
		pat = strings.ReplaceAll(pat, "\\B", "")
		pat = strings.ReplaceAll(pat, "\\a", "")
		pat = strings.ReplaceAll(pat, "\\b", "")
		pat = strings.ReplaceAll(pat, "\\z", "")
		pat = strings.ReplaceAll(pat, "+?", "+")
		pat = strings.ReplaceAll(pat, "*?", "*")
		pat = strings.ReplaceAll(pat, "??", "?")
		pat = strings.ReplaceAll(pat, "}?", "}")
		pat = strings.ReplaceAll(pat, "**", "*")
		pat = strings.ReplaceAll(pat, "+*", "+")
		if pat == pattern {
			return pat
		}
		pattern = pat
	}
}
