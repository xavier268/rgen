package rgen

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xavier268/rgen"
)

func TestGenerate_long(t *testing.T) {
	// using a very complex and large pattern, to test the context timout
	err := dc("a*(c|de|fgh)*", 50)
	if err == nil {
		t.Error("context deadline expected, but did not occur")
	}
}

func ExampleGenerate_ex1() {

	dc("a*", 10)

	// Unordered output:
	// ""
	// "a"
	// "aa"
	// "aaa"
	// "aaaa"
	// "aaaaa"
	// "aaaaaa"
	// "aaaaaaa"
	// "aaaaaaaa"
	// "aaaaaaaaa"
	// "aaaaaaaaaa"
	// <nil>
}

// ==========================================================================================

func dc(pat string, max int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/2)
	defer cancel()

	ch := make(chan string, 10)
	wg := new(sync.WaitGroup)

	// launch async display, until cahnnel is closed or context timeout.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case s, ok := <-ch:
				if !ok {
					return
				}
				fmt.Printf("%q\n", s)
				// iterate ...
			}
		}
	}()

	// launch generation
	err := rgen.Generate(ctx, ch, pat, max)
	fmt.Println(err)

	// signal end of generation to display routine
	close(ch)

	// wait for display to finish printing
	wg.Wait()

	// return nil, unless context was canceled
	return ctx.Err()

}
