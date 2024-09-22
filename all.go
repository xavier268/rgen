package rgen

import (
	"context"
	"iter"

	"github.com/xavier268/rgen/internal/generator"
)

// All returns an iterator over all available strings matching the provided regexp pattern,
// whose lenghth is less or equal to the specified maxlen length.
// No deduplication is performed, if multiple strings can be generated from different path with the same pattern.
// Panic if pattern is not a valid regexp.
func All(pattern string, maxlen int) iter.Seq[string] {

	g, err := generator.NewGenerator(context.Background(), pattern, maxlen)
	if err != nil {
		panic(err)
	}

	return func(yield func(string) bool) {
		for i := 0; i <= maxlen; i++ { // explore all lengths, starting with shorter
			if err := g.Reset(i); err != nil {
				panic(err)
			}
			// generate strings for this length

			for {
				if err := g.Next(); err != nil {
					break
				}
				if !yield(g.Last()) {
					return
				}
			}
		}
	}
}

// Same as All, but the iterator will only generate strings with the provided exact length
func AllExact(pattern string, exactlen int) iter.Seq[string] {

	g, err := generator.NewGenerator(context.Background(), pattern, exactlen)
	if err != nil {
		panic(err)
	}

	return func(yield func(string) bool) {
		if err := g.Reset(exactlen); err != nil {
			panic(err)
		}
		for {
			if err := g.Next(); err != nil {
				return
			}
			if !yield(g.Last()) {
				return
			}
		}
	}
}
