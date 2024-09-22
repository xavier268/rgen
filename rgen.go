package rgen

import (
	"context"
	"iter"
	"sync"

	"github.com/xavier268/rgen/internal/generator"
)

// ================ iterator API ========================

// All returns an iterator over all available strings matching the provided regexp pattern,
// whose lenghth is less or equal to the specified maxlen length.
// No deduplication is performed, if multiple strings can be generated from different path with the same pattern.
// Panic if pattern is not a valid regexp.
// This API is the prefered API.
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

// =========  async API =============================================

// Generate asynchroneously strings up to max length (included) and send them to the channel.
// It is the callers responsibility to read from, and close, the channel.
// The context error is returned, if context was canceled.
// Benchmarks are showing a somewhat limited benefit. Do not choose this API for performance reason, but rather
// to be able to generate a mixture of short and long strings from a complex regexp, until you decide to stop
// using the context.
func Generate(ctx context.Context, ch chan<- string, pattern string, max int) error {
	wg := new(sync.WaitGroup)
	for i := 0; i <= max; i++ { // launch a generator for each length, in parallele.
		wg.Add(1)
		go func() {
			defer wg.Done()
			gen, err := generator.NewGenerator(ctx, pattern, i)
			if err != nil {
				return
			}
			err = gen.Reset(i)
			if err != nil {
				return
			}
			for {
				err = gen.Next()
				if err != nil {
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
					ch <- gen.Last()
				}
			}
		}()
	}
	wg.Wait()
	return ctx.Err()
}
