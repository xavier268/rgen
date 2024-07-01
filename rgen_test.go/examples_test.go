package rgen_test

import (
	"context"
	"fmt"

	"github.com/xavier268/revregex/rgen"
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
