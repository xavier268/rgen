package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/xavier268/revregex"
)

func TestVisual(t *testing.T) {

	revregex.DEBUG = false

	patt := "(ab)(xy)|cde|f(k)|(g|zt)h"

	fmt.Printf("Testing for pattern : %s\n", patt)

	for i := 0; i < 5; i++ {
		g, err := revregex.NewGenerator(context.Background(), patt, i)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println("----", i, "----")
		//"----", i, "----")
		//t.Log(g.String())"
		for {
			f, err := g.Next()
			if err != nil {
				break
			}
			fmt.Println("\t\t\t(", i, ")  -->", f)
		}
	}

}
