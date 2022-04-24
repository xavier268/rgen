package revregex

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestGenStringVisual(t *testing.T) {

	tt := []string{
		"x?",
		"x{3,4}",
		"xxxx?",
		"[c-f]",
		"[^c-f]",
		"[[:alpha:]]",
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
		g := New(rs)
		fmt.Println(g)
	}

}

func TestVerify(t *testing.T) {
	g := New("a*b|c")
	err := g.Verify("c")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenNext(t *testing.T) {

	pats := []struct {
		src string
		cnt int // expected count, negative means do not check.
	}{
		{"a", 1},
		{"abc", 1},
		{"", 1},
		{"ab?", 2},
		{"(ab)c", 1},
		{"(ab)?c", 2},
		{"(ab)c?", 2},
		{"(a)|(b)", 2},
		{"a|(b)", 2},
		{"a|b", 2},
		{"[ab]", 2},
		{"[abc]", 3},
		{"[a-c]", 3},
		{"[i-m]", 5},
		{"[a-ci-m]", 8},
		{"[[:digit:]]", 10},
		{"[[:punct:]]", 32},
		{"[^a]", -1},
		{".", -1},
	}
	it := NewRandInter()
	it.(*rand.Rand).Seed(4242) // fixed seed for reproductibility

	const loop = 1000
	for _, p := range pats {
		s := p.src
		m := make(map[string]int, loop)
		g := New(s)
		fmt.Println(g)
		for i := 0; i < loop; i++ {
			ss := g.Next(it)
			m[ss] += 1
			//fmt.Println(s, "\t->\t", ss)
			err := g.Verify(ss)
			if err != nil {
				fmt.Println(g)
				fmt.Printf("Generated : %#q\n", ss)
				t.Fatal(err)
			}
		}
		i := 0
		for k, v := range m {
			i++
			if i < 20 {
				fmt.Printf("%3.1f%%\t%s\t%#v\n", 100*float64(v)/loop, k, []byte(k))
			} else {
				fmt.Printf("[....]  truncating a total of %d values\n", len(m))
				break
			}
		}
		if p.cnt >= 0 && len(m) != p.cnt {
			t.Fatalf("the total number of different strings (%d) does not match expectation (%d)", len(m), p.cnt)
		}
		fmt.Println()
	}
}
