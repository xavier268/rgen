package revregex

import (
	"fmt"
	"testing"
)

func TestDump(t *testing.T) {

	tt := []string{
		"[c-f]",
		"[^c-f]",
		"abcdefg",
		"^abcdefg",
		"abcdefg$",
		"\\+",
		"a|b",
		"a|b|c",
		"a|",
		"(a|)",
		"a(tv)+f*|6+xx[a-f]z",
	}

	for i, rs := range tt {
		// Display tree
		fmt.Print(i, "\t")
		Dump(rs)
		fmt.Println()
	}

}
