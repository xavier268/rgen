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

func ExampleNewGenerator_alternate4() {
	do("[0-3]|[7-8]", 4)

	// Output:
	// Testing for pattern : "[0-3]|[7-8]"
	// (1) --> "0"
	// (1) --> "1"
	// (1) --> "2"
	// (1) --> "3"
	// (1) --> "7"
	// (1) --> "8"
}

func ExampleNewGenerator_concat1() {
	do("(g|zt)h", 5)
	// Unordered output:
	// Testing for pattern : "(g|zt)h"
	// (2) --> "gh"
	// (3) --> "zth"

}

func ExampleNewGenerator_quest1() {

	do("a?", 4)
	// Output:
	// Testing for pattern : "a?"
	// (0) --> ""
	// (1) --> "a"
}
func ExampleNewGenerator_quest2() {

	do("a?b", 4)
	// Output:
	// Testing for pattern : "a?b"
	// (1) --> "b"
	// (2) --> "ab"
}
func ExampleNewGenerator_quest2b() {

	do("ba?", 4)
	// Output:
	// Testing for pattern : "ba?"
	// (1) --> "b"
	// (2) --> "ba"
}

func ExampleNewGenerator_quest3() {
	// revregex.DEBUG = true
	do("ab?c", 6)

	// Output:
	// Testing for pattern : "ab?c"
	// (2) --> "ac"
	// (3) --> "abc"
}

func ExampleNewGenerator_quest4() {
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
	// 	Testing for pattern : "[0-3]*"
	// (0) --> ""
	// (1) --> "0"
	// (1) --> "1"
	// (1) --> "2"
	// (1) --> "3"
	// (2) --> "00"
	// (2) --> "01"
	// (2) --> "02"
	// (2) --> "03"
	// (2) --> "10"
	// (2) --> "11"
	// (2) --> "12"
	// (2) --> "13"
	// (2) --> "20"
	// (2) --> "21"
	// (2) --> "22"
	// (2) --> "23"
	// (2) --> "30"
	// (2) --> "31"
	// (2) --> "32"
	// (2) --> "33"
	// (3) --> "000"
	// (3) --> "001"
	// (3) --> "002"
	// (3) --> "003"
	// (3) --> "010"
	// (3) --> "011"
	// (3) --> "012"
	// (3) --> "013"
	// (3) --> "020"
	// (3) --> "021"
	// (3) --> "022"
	// (3) --> "023"
	// (3) --> "030"
	// (3) --> "031"
	// (3) --> "032"
	// (3) --> "033"
	// (3) --> "100"
	// (3) --> "101"
	// (3) --> "102"
	// (3) --> "103"
	// (3) --> "110"
	// (3) --> "111"
	// (3) --> "112"
	// (3) --> "113"
	// (3) --> "120"
	// (3) --> "121"
	// (3) --> "122"
	// (3) --> "123"
	// (3) --> "130"
	// (3) --> "131"
	// (3) --> "132"
	// (3) --> "133"
	// (3) --> "200"
	// (3) --> "201"
	// (3) --> "202"
	// (3) --> "203"
	// (3) --> "210"
	// (3) --> "211"
	// (3) --> "212"
	// (3) --> "213"
	// (3) --> "220"
	// (3) --> "221"
	// (3) --> "222"
	// (3) --> "223"
	// (3) --> "230"
	// (3) --> "231"
	// (3) --> "232"
	// (3) --> "233"
	// (3) --> "300"
	// (3) --> "301"
	// (3) --> "302"
	// (3) --> "303"
	// (3) --> "310"
	// (3) --> "311"
	// (3) --> "312"
	// (3) --> "313"
	// (3) --> "320"
	// (3) --> "321"
	// (3) --> "322"
	// (3) --> "323"
	// (3) --> "330"
	// (3) --> "331"
	// (3) --> "332"
	// (3) --> "333"
}
func ExampleNewGenerator_class4() {
	do("[0-2]*[7-9]*", 2)

	// Output:
	// Testing for pattern : "[0-2]*[7-9]*"
	// (0) --> ""
	// (1) --> "0"
	// (1) --> "1"
	// (1) --> "2"
	// (1) --> "7"
	// (1) --> "8"
	// (1) --> "9"

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
