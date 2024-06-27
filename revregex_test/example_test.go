package revregex_test

import (
	"context"
	"fmt"

	"github.com/xavier268/revregex"
)

func ExampleNewGenerator_alternate_ex1() {

	do("(ab)|cde|f|gh", 6)

	// Output:
	// Testing for pattern : (ab)|cde|f|gh
	// (1) --> "f"
	// (2) --> "ab"
	// (2) --> "gh"
	// (3) --> "cde"
}

func ExampleNewGenerator_concat_ex1() {

	do("(ab)(xy)|cde|f(k)|(g|zt)h", 6)

	// Output:
	// Testing for pattern : (ab)(xy)|cde|f(k)|(g|zt)h
	// (2) --> "fk"
	// (2) --> "gh"
	// (3) --> "cde"
	// (3) --> "zth"
	// (4) --> "abxy"
}

func ExampleNewGenerator_quest_ex1() {
	// revregex.DEBUG = true
	do("ab?c", 6)

	// Output:
	// Testing for pattern : ab?c
	// (2) --> "ac"
	// (3) --> "abc"
}

func ExampleNewGenerator_quest_ex2() {
	// revregex.DEBUG = true
	do("a(b|ut)?c", 6)

	// Output:
	// Testing for pattern : a(b|ut)?c
	// (2) --> "ac"
	// (3) --> "abc"
	// (4) --> "autc"
}

func ExampleNewGenerator_star_ex1() {

	do("ab*c", 6)

	// Output:
	// Testing for pattern : ab*c
	// (2) --> "ac"
	// (3) --> "abc"
	// (4) --> "abbc"
	// (5) --> "abbbc"
}

func ExampleNewGenerator_star_ex2() {

	do("a*b*c", 6)

	// Output:
	// 	Testing for pattern : a*b*c
	// (1) --> "c"
	// (2) --> "ac"
	// (2) --> "bc"
	// (3) --> "aac"
	// (3) --> "abc"
	// (3) --> "bbc"
	// (4) --> "aaac"
	// (4) --> "aabc"
	// (4) --> "abbc"
	// (4) --> "bbbc"
	// (5) --> "aaaac"
	// (5) --> "aaabc"
	// (5) --> "aabbc"
	// (5) --> "abbbc"
	// (5) --> "bbbbc"
}

func ExampleNewGenerator_star_ex3() {

	do("a*b*", 3)

	// Output:
	// Testing for pattern : a*b*
	// (0) --> ""
	// (1) --> "a"
	// (1) --> "b"
	// (2) --> "aa"
	// (2) --> "ab"
	// (2) --> "bb"

}

// =====================
func do(patt string, n int) {
	fmt.Printf("Testing for pattern : %s\n", patt)

	for i := 0; i < n; i++ {
		g, err := revregex.NewGenerator(context.Background(), patt, i)
		if err != nil {
			panic(err)
		}

		for {
			f, err := g.Next()
			if err != nil {
				break
			}
			fmt.Printf("(%d) --> %q\n", i, f)
		}
	}
}
