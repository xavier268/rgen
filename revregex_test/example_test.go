package revregex_test

import (
	"context"
	"fmt"
	"regexp/syntax"

	"github.com/xavier268/revregex"
)

func ExampleNewGenerator_alternate1() {

	do("(ab)|cde|f|gh", 6)

	// Unordered output:
	// Testing for pattern : "(ab)|cde|f|gh"
	// (1) --> "f"
	// (2) --> "ab"
	// (2) --> "gh"
	// (3) --> "cde"
}

func ExampleNewGenerator_alternate2() {
	// revregex.DEBUG = true
	do("(a|b|c)(e|f|g)", 3)
	// Unordered output:
	// Testing for pattern : "(a|b|c)(e|f|g)"
	// (2) --> "ae"
	// (2) --> "af"
	// (2) --> "ag"
	// (2) --> "be"
	// (2) --> "bf"
	// (2) --> "bg"
	// (2) --> "ce"
	// (2) --> "cf"
	// (2) --> "cg"
}

func ExampleNewGenerator_alternate3() {
	// revregex.DEBUG = true
	do("(a|bc|d|efg)", 5)
	// Unordered output:
	// Testing for pattern : "(a|bc|d|efg)"
	// (1) --> "a"
	// (1) --> "d"
	// (2) --> "bc"
	// (3) --> "efg"
}

func ExampleNewGenerator_concat1() {
	do("(g|zt)h", 5)
	// Unordered output:
	// Testing for pattern : "(g|zt)h"
	// (2) --> "gh"
	// (3) --> "zth"

}
func ExampleNewGenerator_concat2() {

	do("(ab)(xy)|cde|f(k)|(g|zt)h", 6)

	// Unordered output:
	// Testing for pattern : "(ab)(xy)|cde|f(k)|(g|zt)h"
	// (2) --> "fk"
	// (2) --> "gh"
	// (3) --> "cde"
	// (3) --> "zth"
	// (4) --> "abxy"
}

func ExampleNewGenerator_quest1() {
	// revregex.DEBUG = true
	do("ab?c", 6)

	// Output:
	// Testing for pattern : "ab?c"
	// (2) --> "ac"
	// (3) --> "abc"
}

func ExampleNewGenerator_quest2() {
	// revregex.DEBUG = true
	do("a(b|ut)?c", 6)

	// Output:
	// Testing for pattern : "a(b|ut)?c"
	// (2) --> "ac"
	// (3) --> "abc"
	// (4) --> "autc"
}

func ExampleNewGenerator_star1() {

	do("ab*c", 6)

	// Output:
	// Testing for pattern : "ab*c"
	// (2) --> "ac"
	// (3) --> "abc"
	// (4) --> "abbc"
	// (5) --> "abbbc"
}

func ExampleNewGenerator_star2() {

	do("a*b*c", 6)

	// Output:
	// 	Testing for pattern : "a*b*c"
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

func ExampleNewGenerator_star3() {

	do("a*b*", 3)

	// Output:
	// Testing for pattern : "a*b*"
	// (0) --> ""
	// (1) --> "a"
	// (1) --> "b"
	// (2) --> "aa"
	// (2) --> "ab"
	// (2) --> "bb"

}

func ExampleNewGenerator_star4() {
	do("(a|b)*", 3)

	// Output:
	// Testing for pattern : "(a|b)*"
	// (0) --> ""
	// (1) --> "a"
	// (1) --> "b"
	// (2) --> "aa"
	// (2) --> "ab"
	// (2) --> "ba"
}

func ExampleNewGenerator_limitedRange() {
	do("a{1,4}", 6)

	// Output:
	// Testing for pattern : "a{1,4}"
	// (1) --> "a"
	// (2) --> "aa"
	// (3) --> "aaa"
	// (4) --> "aaaa"
}

func ExampleNewGenerator_plus1() {
	do("a+", 4)

	// Output:
	// Testing for pattern : "a+"
	// (1) --> "a"
	// (2) --> "aa"
	// (3) --> "aaa"
}

func ExampleNewGenerator_plus2() {
	do("(a+)x(b*)", 5)

	// Output:
	// Testing for pattern : "(a+)x(b*)"
	// (2) --> "ax"
	// (3) --> "aax"
	// (3) --> "axb"
	// (4) --> "aaax"
	// (4) --> "aaxb"
	// (4) --> "axbb"
}
func ExampleNewGenerator_class1() {
	do("[a-c]", 5)
	do("[v-v]", 5)
	do("[xyz]", 5)
	do("[z]", 5)

	// Output:
	// Testing for pattern : "[a-c]"
	// (1) --> "a"
	// (1) --> "b"
	// (1) --> "c"
	// Testing for pattern : "[v-v]"
	// (1) --> "v"
	// Testing for pattern : "[xyz]"
	// (1) --> "x"
	// (1) --> "y"
	// (1) --> "z"
	// Testing for pattern : "[z]"
	// (1) --> "z"
}

func ExampleNewGenerator_class2() {
	do("[a-ct-vx]", 5)
	do("[\n\r\t]", 5)

	// Output:
	// Testing for pattern : "[a-ct-vx]"
	// (1) --> "a"
	// (1) --> "b"
	// (1) --> "c"
	// (1) --> "t"
	// (1) --> "u"
	// (1) --> "v"
	// (1) --> "x"
	// Testing for pattern : "[\n\r\t]"
	// (1) --> "\t"
	// (1) --> "\n"
	// (1) --> "\r"
}
func ExampleNewGenerator_class3() {
	do("[0-3]*", 4)
	// Output:
}
func ExampleNewGenerator_class4() {
	do("http(s)?://[0-9]{1,3}|(ab*a)", 15)

	// Output:
	// Testing for pattern : "http(s)?://[0-9]{1,3}|(ab*a)"
	// (2) --> "aa"
	// (3) --> "aba"
	// (4) --> "abba"
	// (5) --> "abbba"
	// (6) --> "abbbba"
	// (7) --> "abbbbba"
	// (8) --> "http://0"
	// (8) --> "abbbbbba"
	// (9) --> "https://0"
	// (9) --> "http://00"
	// (9) --> "abbbbbbba"
	// (10) --> "https://00"
	// (10) --> "http://000"
	// (10) --> "http://001"
	// (10) --> "http://002"
	// (10) --> "http://003"
	// (10) --> "http://004"
	// (10) --> "http://005"
	// (10) --> "http://006"
	// (10) --> "http://007"
	// (10) --> "http://008"
	// (10) --> "http://009"
	// (10) --> "abbbbbbbba"
	// (11) --> "https://000"
	// (11) --> "https://001"
	// (11) --> "https://002"
	// (11) --> "https://003"
	// (11) --> "https://004"
	// (11) --> "https://005"
	// (11) --> "https://006"
	// (11) --> "https://007"
	// (11) --> "https://008"
	// (11) --> "https://009"
	// (11) --> "abbbbbbbbba"
	// (12) --> "abbbbbbbbbba"
	// (13) --> "abbbbbbbbbbba"
	// (14) --> "abbbbbbbbbbbba"
}

func ExampleNewGenerator_simplified_parsing() {
	pats := []string{"a{1,5}", "a+", "a*", "a+a+", "a**", "a++"}
	for _, s := range pats {
		re, err := syntax.Parse(s, syntax.POSIX)
		if err != nil {
			panic(err)
		}
		res := re.Simplify()
		fmt.Printf("%q --parse--> %q --simplify--> %q\n", s, re, res)
	}

	// Output:
	// "a{1,5}" --parse--> "a{1,5}" --simplify--> "a(?:a(?:a(?:aa?)?)?)?"
	// "a+" --parse--> "a+" --simplify--> "a+"
	// "a*" --parse--> "a*" --simplify--> "a*"
	// "a+a+" --parse--> "a+a+" --simplify--> "a+a+"
	// "a**" --parse--> "(?:a*)*" --simplify--> "a*"
	// "a++" --parse--> "(?:a+)+" --simplify--> "a+"
}

// =====================
func do(patt string, n int) {
	fmt.Printf("Testing for pattern : %q\n", patt)

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
