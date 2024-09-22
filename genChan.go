package rgen

import (
	"context"
	"sync"

	"github.com/xavier268/rgen/internal/generator"
)

// Generate asynchroneously strings up to max length (included) and send them to the channel.
// It is the callers responsibility to read from, and close, the channel.
// The context error is returned, if context was canceled.
func Generate(ctx context.Context, ch chan<- string, pattern string, max int) error {
	wg := new(sync.WaitGroup)
	for i := 0; i <= max; i++ {
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
