package rgen_test

import (
	"context"
	"fmt"

	"github.com/xavier268/rgen"
)

func ExampleNewGenerator_alternate1() {

	do("a|bc", 3)

	// Unordered output:
	// Testing for pattern : "a|bc"
	// (1) --> "a"
	// (2) --> "bc"
}

func ExampleNewGenerator_alternate2() {

	do("a|bc|cde|ef|g", 3)

	// Unordered output:
	// Testing for pattern : "a|bc|cde|ef|g"
	// (1) --> "a"
	// (1) --> "g"
	// (2) --> "bc"
	// (2) --> "ef"
	// (3) --> "cde"
}
func ExampleNewGenerator_alternate3() {
	do("a|((bc|cde)|ef)|g", 3)

	// Unordered output:
	// Testing for pattern : "a|((bc|cde)|ef)|g"
	// (1) --> "a"
	// (1) --> "g"
	// (2) --> "bc"
	// (2) --> "ef"
	// (3) --> "cde"
}

func ExampleNewGenerator_class1() {
	do("[a-f]", 3)

	// Unordered output:
	// Testing for pattern : "[a-f]"
	// (1) --> "a"
	// (1) --> "b"
	// (1) --> "c"
	// (1) --> "d"
	// (1) --> "e"
	// (1) --> "f"
}

func ExampleNewGenerator_class2() {
	do("[ak]", 3)

	// Unordered output:
	// Testing for pattern : "[ak]"
	// (1) --> "a"
	// (1) --> "k"
}

func ExampleNewGenerator_class3() {
	do("[a]", 3)

	// Unordered output:
	// Testing for pattern : "[a]"
	// (1) --> "a"
}

func ExampleNewGenerator_class4() {
	do("[a-bt-vz]", 3)

	// Unordered output:
	// Testing for pattern : "[a-bt-vz]"
	// (1) --> "a"
	// (1) --> "b"
	// (1) --> "t"
	// (1) --> "u"
	// (1) --> "v"
	// (1) --> "z"
}

func ExampleNewGenerator_quest1() {
	do("a?", 3)

	// Unordered output:
	// Testing for pattern : "a?"
	// (0) --> ""
	// (1) --> "a"
}

func ExampleNewGenerator_quest2() {
	do("a?|b", 3)

	// Unordered output:
	// Testing for pattern : "a?|b"
	// (0) --> ""
	// (1) --> "a"
	// (1) --> "b"
}

func ExampleNewGenerator_quest3() {
	do("[a-b]?", 3)

	// Unordered output:
	// Testing for pattern : "[a-b]?"
	// (0) --> ""
	// (1) --> "a"
	// (1) --> "b"
}

func ExampleNewGenerator_concat1() {
	do("[a-b]c", 3)

	// Unordered output:
	// Testing for pattern : "[a-b]c"
	// (2) --> "ac"
	// (2) --> "bc"

}
func ExampleNewGenerator_concat2() {
	do("[a-b][c-d]", 3)

	// Unordered output:
	// Testing for pattern : "[a-b][c-d]"
	// (2) --> "ac"
	// (2) --> "bc"
	// (2) --> "ad"
	// (2) --> "bd"
}

func ExampleNewGenerator_concat3() {
	do("a?c", 3)

	// Unordered output:
	// Testing for pattern : "a?c"
	// (1) --> "c"
	// (2) --> "ac"

}
func ExampleNewGenerator_concat4() {
	do("(a|b)c", 3)

	// Unordered output:
	// Testing for pattern : "(a|b)c"
	// (2) --> "ac"
	// (2) --> "bc"

}

func ExampleNewGenerator_concat5() {
	do("(a|bz)c", 3)

	// Unordered output:
	// Testing for pattern : "(a|bz)c"
	// (2) --> "ac"
	// (3) --> "bzc"

}

func ExampleNewGenerator_concat6() {
	do("[0-3]?[a-b]?[x-y]?", 5)

	// Unordered output:
	// Testing for pattern : "[0-3]?[a-b]?[x-y]?"
	// (0) --> ""
	// (1) --> "x"
	// (1) --> "y"
	// (1) --> "0"
	// (1) --> "1"
	// (1) --> "2"
	// (1) --> "3"
	// (1) --> "a"
	// (1) --> "b"
	// (2) --> "0x"
	// (2) --> "1x"
	// (2) --> "2x"
	// (2) --> "3x"
	// (2) --> "0y"
	// (2) --> "1y"
	// (2) --> "2y"
	// (2) --> "3y"
	// (2) --> "ax"
	// (2) --> "bx"
	// (2) --> "ay"
	// (2) --> "by"
	// (2) --> "0a"
	// (2) --> "1a"
	// (2) --> "2a"
	// (2) --> "3a"
	// (2) --> "0b"
	// (2) --> "1b"
	// (2) --> "2b"
	// (2) --> "3b"
	// (3) --> "0ax"
	// (3) --> "1ax"
	// (3) --> "2ax"
	// (3) --> "3ax"
	// (3) --> "0bx"
	// (3) --> "1bx"
	// (3) --> "2bx"
	// (3) --> "3bx"
	// (3) --> "0ay"
	// (3) --> "1ay"
	// (3) --> "2ay"
	// (3) --> "3ay"
	// (3) --> "0by"
	// (3) --> "1by"
	// (3) --> "2by"
	// (3) --> "3by"
}

func ExampleNewGenerator_concat7() {
	do("(a|cd|efg)", 5)
	// Unordered output:
	// Testing for pattern : "(a|cd|efg)"
	// (1) --> "a"
	// (2) --> "cd"
	// (3) --> "efg"
}

func ExampleNewGenerator_star1() {
	do("a*", 4)
	// Unordered output:
	// Testing for pattern : "a*"
	// (0) --> ""
	// (1) --> "a"
	// (2) --> "aa"
	// (3) --> "aaa"
	// (4) --> "aaaa"
}

func ExampleNewGenerator_star2() {
	do("a*b*", 4)
	// Unordered output:
	// Testing for pattern : "a*b*"
	// (0) --> ""
	// (1) --> "b"
	// (1) --> "a"
	// (2) --> "bb"
	// (2) --> "ab"
	// (2) --> "aa"
	// (3) --> "bbb"
	// (3) --> "abb"
	// (3) --> "aab"
	// (3) --> "aaa"
	// (4) --> "bbbb"
	// (4) --> "abbb"
	// (4) --> "aabb"
	// (4) --> "aaab"
	// (4) --> "aaaa"
}

func ExampleNewGenerator_star3() {
	do("(a|cd|efg)*", 4)
	// Unordered output:
	// Testing for pattern : "(a|cd|efg)*"
	// (0) --> ""
	// (1) --> "a"
	// (2) --> "cd"
	// (2) --> "aa"
	// (3) --> "efg"
	// (3) --> "cda"
	// (3) --> "acd"
	// (3) --> "aaa"
	// (4) --> "efga"
	// (4) --> "cdcd"
	// (4) --> "cdaa"
	// (4) --> "aefg"
	// (4) --> "acda"
	// (4) --> "aacd"
	// (4) --> "aaaa"
}

func ExampleNewGenerator_plus1() {
	do("a+", 4)
	// Unordered output:
	// Testing for pattern : "a+"
	// (1) --> "a"
	// (2) --> "aa"
	// (3) --> "aaa"
	// (4) --> "aaaa"
}

func ExampleNewGenerator_plus2() {
	do("a+b+", 4)
	// Unordered output:
	// Testing for pattern : "a+b+"
	// (2) --> "ab"
	// (3) --> "abb"
	// (3) --> "aab"
	// (4) --> "abbb"
	// (4) --> "aabb"
	// (4) --> "aaab"
}

//========================================================================

func do(patt string, n int) {
	fmt.Printf("Testing for pattern : %q\n", patt)

	// set max length
	g, err := rgen.NewGenerator(context.Background(), patt, n)
	if err != nil {
		panic(err)
	}

	// try successively all length
	for i := 0; i <= n; i++ {
		if err := g.Reset(i); err != nil {
			panic(err)
		}
		// show all result strings
		for {
			if err := g.Next(); err != nil {
				break
			}
			fmt.Printf("(%d) --> %q\n", i, g.Last())
		}
	}
}
