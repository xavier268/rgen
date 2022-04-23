package revregex

import (
	"fmt"
	"math/big"
	"math/rand"
	"regexp/syntax"
	"strings"
	"time"
)

// Gen can efficiently and deterministically generate strings that match a given regexp.
type Gen struct {
	// source string for the regexp.
	source string
	// root parsed tree
	tree *syntax.Regexp
	// random generator.
	rand *rand.Rand
	// max length when generating with * or + metasymbols. Has no impact on the actual length of the generated strings.
	// set to 0 to disable * and +
	// set to a negative value to have no formal limit, but an exponential law.
	ml int
}

// New creates a new generator.
// It will panic if the regexp provided is not syntacly correct.
// Use POSIX syntax.
func New(source string) *Gen {
	var err error
	g := new(Gen)
	g.source = source
	g.tree, err = syntax.Parse(source, syntax.POSIX)
	if err != nil {
		panic(err)
	}
	g.tree = g.tree.Simplify()
	g.rand = rand.New(rand.NewSource(time.Hour.Milliseconds()))
	g.ml = 6 // see if we want to keep this default value ?
	return g
}

// Next provides a random string matching the regexp.
func (g *Gen) Next() string {
	panic("to do")
}

// NextI provides a determistic string matching the regexp.
// i should be strictly positive.
// The bigint returned it garanteed to be smaller or equal to the initial i value.
func (g *Gen) NextI(i *big.Int) (string, *big.Int, error) {
	panic("to do")
}

var ErrEntropyExhausted = fmt.Errorf("unsufficient entropy available")

func Dump(s string) {

	fmt.Printf("%q\n", s)
	re, err := syntax.Parse(s, 0)
	if err != nil {
		panic(err)
	}

	var b strings.Builder
	dump(&b, re.Simplify()) // use re.Simplify ?
	fmt.Println(b.String())

	fmt.Println("Freedom : ", freedom(re))

}

func dump(b *strings.Builder, re *syntax.Regexp) {
	if re == nil {
		fmt.Fprint(b, nil)
		return
	}
	switch re.Op {

	case
		syntax.OpLiteral, syntax.OpCharClass:
		fmt.Fprintf(b, "%s(%q)", re.Op, re.Rune)

	default:

		fmt.Fprint(b, re.Op, "(")
		for _, rs := range re.Sub {
			dump(b, rs)
		}
		fmt.Fprint(b, ")")
	}
}

// freedom degree of current tree
// TODO - completer/verifier !!
func freedom(re *syntax.Regexp) int64 {

	if re == nil {
		return 0
	}
	switch re.Op {

	case syntax.OpAlternate:
		return int64(len(re.Sub))
	case syntax.OpQuest:
		return 2 * freedom(re.Sub[0])
	case syntax.OpCapture:
		f := int64(1)
		for _, rr := range re.Sub {
			f *= freedom(rr)
		}
		return f
	case syntax.OpConcat:
		f := int64(1)
		for _, rr := range re.Sub {
			f *= freedom(rr)
		}
		return f
	case syntax.OpCharClass:
		var f int64
		for i := 0; i+1 < len(re.Rune); i += 2 {
			f += int64(re.Rune[i+1]) - int64(re.Rune[i]) + 1
		}
		return f
	default:
		return 1
	}
}
