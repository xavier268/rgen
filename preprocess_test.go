package revregex

import (
	"fmt"
	"regexp/syntax"
	"testing"
)

func TestPreProcessVisual(t *testing.T) {

	show("(g|zt)h")
	show("a(b)(cd)k*(efg)u{3,5}c+")

}

func show(pat string) {
	fmt.Println("Analysing : ", pat)
	fmt.Println(pat)
	fmt.Println("Parsing ...")
	re, err := syntax.Parse(pat, syntax.POSIX)
	if err != nil {
		panic(err)
	}
	fmt.Println(re)
	fmt.Println("Simplifying ...")
	re = re.Simplify()
	fmt.Println(re)
	fmt.Println(Dump(re))
	fmt.Println("Preprocessing ...")
	re = preProcess(re)
	fmt.Println(re)
	fmt.Println(Dump(re))
}
