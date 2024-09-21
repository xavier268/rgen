package rgen_test

import (
	"fmt"

	"github.com/xavier268/rgen"
)

func ExampleAll() {

	pattern := "a?(b|c)*"
	maxlen := 4

	fmt.Println("Pattern:", pattern, "Max length: up to", maxlen, "(included)")

	// Iterate over all strings generated from pattern up to maxlen length ...
	for s := range rgen.All(pattern, maxlen) {
		fmt.Println(s)
	}

	// Output:
	// Pattern: a?(b|c)* Max length: up to 4 (included)
	//
	// b
	// c
	// a
	// bb
	// cb
	// bc
	// cc
	// ab
	// ac
	// bbb
	// cbb
	// bcb
	// ccb
	// bbc
	// cbc
	// bcc
	// ccc
	// abb
	// acb
	// abc
	// acc
	// bbbb
	// cbbb
	// bcbb
	// ccbb
	// bbcb
	// cbcb
	// bccb
	// cccb
	// bbbc
	// cbbc
	// bcbc
	// ccbc
	// bbcc
	// cbcc
	// bccc
	// cccc
	// abbb
	// acbb
	// abcb
	// accb
	// abbc
	// acbc
	// abcc
	// accc

}

func ExampleAllExact() {

	pattern := "a?(b|c)*"
	exactlen := 4

	fmt.Println("Pattern:", pattern, "Exact length:", exactlen)

	// Iterate over all strings generated from pattern up to maxlen length ...
	for s := range rgen.AllExact(pattern, exactlen) {
		fmt.Println(s)
	}

	// Note how the empty string is not generated anay more below ...

	// Output:
	// Pattern: a?(b|c)* Exact length: 4
	// bbbb
	// cbbb
	// bcbb
	// ccbb
	// bbcb
	// cbcb
	// bccb
	// cccb
	// bbbc
	// cbbc
	// bcbc
	// ccbc
	// bbcc
	// cbcc
	// bccc
	// cccc
	// abbb
	// acbb
	// abcb
	// accb
	// abbc
	// acbc
	// abcc
	// accc

}

func ExampleDedup() {

	pattern := "a+b?a+"
	len := 4
	fmt.Println("Pattern:", pattern, "Exact length:", len)

	// Iterate over all strings generated from pattern up to maxlen length ...
	// deduplicating using the provided map deduper.
	// for s := range rgen.Dedup(rgen.AllExact(pattern, len), rgen.NewDedupMap()) {
	fmt.Println("\nWithout Dedup")
	for s := range rgen.AllExact(pattern, len) {
		fmt.Println(s)
	}
	fmt.Println("\nWith Dedup")
	for s := range rgen.Dedup(rgen.AllExact(pattern, len), rgen.NewDedupMap()) {
		fmt.Println(s)
	}

	// Output:
	// Pattern: a+b?a+ Exact length: 4
	//
	// Without Dedup
	// aaaa
	// aaaa
	// aaaa
	// abaa
	// aaba
	//
	// With Dedup
	// aaaa
	// abaa
	// aaba
}