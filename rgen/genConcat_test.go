package rgen

import (
	"fmt"
)

func Example_incSplitConcat() {

	lls := [][]int{
		{0, 0, 4},
		{0, 3},
		{2},
		{},
		{0, 0, 0, 1},
	}

	for _, ll := range lls {
		n := sum(ll)
		var err error
		fmt.Println()
		fmt.Println(ll) // initial value
		for err == nil {
			err = incSplitConcat(ll)
			if n != sum(ll) {
				panic("incSplitConcat changed the sum of the split")
			}
			fmt.Println(ll, err)
		}
	}

	// Output:
	//
	// [0 0 4]
	// [1 0 3] <nil>
	// [2 0 2] <nil>
	// [3 0 1] <nil>
	// [4 0 0] <nil>
	// [0 1 3] <nil>
	// [1 1 2] <nil>
	// [2 1 1] <nil>
	// [3 1 0] <nil>
	// [0 2 2] <nil>
	// [1 2 1] <nil>
	// [2 2 0] <nil>
	// [0 3 1] <nil>
	// [1 3 0] <nil>
	// [0 4 0] <nil>
	// [0 0 4] done
	//
	// [0 3]
	// [1 2] <nil>
	// [2 1] <nil>
	// [3 0] <nil>
	// [0 3] done
	//
	// [2]
	// [2] done
	//
	// []
	// [] done
	//
	// [0 0 0 1]
	// [1 0 0 0] <nil>
	// [0 1 0 0] <nil>
	// [0 0 1 0] <nil>
	// [0 0 0 1] done

}
