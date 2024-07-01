package rgen

import "fmt"

func Example_incSplitStar() {
	ss := [][]int{
		{4},
		{},
		{1},
	}
	for _, s := range ss {
		var err error
		fmt.Println("\nStarting with ", s)
		for err == nil {
			s, err = incSplitStar(s)
			fmt.Println(s, err)
		}
	}

	// Output:
	// Starting with  [4]
	// [3 1] <nil>
	// [2 2] <nil>
	// [2 1 1] <nil>
	// [1 3] <nil>
	// [1 2 1] <nil>
	// [1 1 2] <nil>
	// [1 1 1 1] <nil>
	// [] done
	//
	// Starting with  []
	// [] done
	//
	// Starting with  [1]
	// [] done
}
