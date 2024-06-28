package revregex_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xavier268/revregex"
)

func TestMain1(t *testing.T) {

	// set timeout to 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// create a generator for the given length.
	g, err := revregex.NewGenerator(ctx, `http(s?)://(www\.)?[0-9]+`, 20)
	if err != nil {
		panic("cannot create a generator")
	}

	// display results ...
	res := ""
	for err == nil {
		res, err = g.Next()
		fmt.Println(res)
	}

	// Report error, informing about timeout
	fmt.Println(err, ", ", ctx.Err())

	if ctx.Err() == nil {
		t.Fatalf("Timeout should have occured")
	}

}

func TestMain2(t *testing.T) {

	// set timeout to 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create channel
	out := make(chan string)

	var wg sync.WaitGroup

	// start a go routine to generate results ...
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting to generate")
		// create a generator for the given length.
		err := revregex.Generate(ctx, `[0-9]*`, 4, out)
		if err != nil {
			fmt.Println("Error: ", err, ctx.Err())
		}
		close(out) // close channel when done.  // important
	}()

	// start a go routine to consume results ...
	wg.Add(1)
	go func() {
		fmt.Println("starting to consume")
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context done - stop reading")
			case res, ok := <-out:
				if ok {
					fmt.Println(res)
				} else {
					fmt.Println("Channel closed")
					return // important !
				}
			}
		}
	}()

	// wait for both go routines to complete()  {
	wg.Wait()

	if ctx.Err() != nil {
		t.Fatalf("should not time out : %v", ctx.Err())
	}
}
