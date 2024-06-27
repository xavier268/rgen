package revregex_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/xavier268/revregex"
)

func TestMain(t *testing.T) {

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

}
