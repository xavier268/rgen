package revregex

import (
	"fmt"
	"regexp/syntax"
	"strings"
)

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
