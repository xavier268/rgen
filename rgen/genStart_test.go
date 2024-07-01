package rgen

import "fmt"

func Example_incSplitStar() {
	s := []int{4}
	var err error
	var v []int

	for err == nil {
		v, err = incSplitStar(s)
		fmt.Println(v, err)
	}

	// Output:

}
