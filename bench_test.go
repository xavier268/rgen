package revregex

import (
	"testing"
)

// 2022-04-24 : go test --bench=. --benchmem --cover
//
// goos: linux
// goarch: amd64
// pkg: github.com/xavier268/revregex
// cpu: Intel(R) Core(TM) i7-10700 CPU @ 2.90GHz
// BenchmarkNextBoundedRaw-8                1000000              1008 ns/op             231 B/op         20 allocs/op
// BenchmarkNextBoundedSimplified-8         1658962               672.1 ns/op           147 B/op         12 allocs/op
// BenchmarkNextUnboundedRaw-8              1921507               628.5 ns/op           138 B/op         11 allocs/op
// BenchmarkNextUnboundedSimplified-8       1586078               639.0 ns/op           138 B/op         11 allocs/op
// PASS
// coverage: 81.1% of statements

var result string // prevents compiler overoptimization ;-)

func BenchmarkNextBoundedRaw(b *testing.B) {
	s := "a*(b|c?d{2,5})e{1,10}"
	//fmt.Printf("Benchmarking %#q\n", s)
	g := NewGen(s)
	it := NewRandChooserSeed(4242)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = g.Next(it)
	}

}

func BenchmarkNextBoundedSimplified(b *testing.B) {
	s := "a*(b|c?d{2,5})e{1,10}"
	//fmt.Printf("Benchmarking %#q\n", s)
	g := NewGenSimpl(s)
	it := NewRandChooserSeed(4242)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = g.Next(it)
	}

}

func BenchmarkNextUnboundedRaw(b *testing.B) {
	s := "a*(b|c?d+)e+"
	//fmt.Printf("Benchmarking %#q\n", s)
	g := NewGen(s)
	it := NewRandChooserSeed(4242)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = g.Next(it)
	}
}

func BenchmarkNextUnboundedSimplified(b *testing.B) {
	s := "a*(b|c?d+)e+"
	//fmt.Printf("Benchmarking %#q\n", s)
	g := NewGenSimpl(s)
	it := NewRandChooserSeed(4242)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = g.Next(it)
	}
}
