package revregex

import (
	"fmt"
	"testing"
)

// go test -fuzz=.

func FuzzNextGenBytes(f *testing.F) {

	tb := [][]byte{
		//{},
		{255, 34, 18},
		{0, 67},
	}
	for _, tt := range tb {
		f.Add(tt)
	}

	pat := "(a|b+|c)(dc?){2,5}"

	ft := func(t *testing.T, bb []byte) {
		g, err := NewGen(pat)
		if err != nil {
			panic(err)
		}
		it := NewBytesChooser(bb)
		r1 := g.Next(it)
		err = g.Verify(r1)
		if err != nil {
			t.Errorf("String <%#q> did not verify : %v", r1, err)
		}
		r2 := g.Next(it)
		err = g.Verify(r2)
		if err != nil {
			t.Errorf("String <%#q> did not verify : %v", r2, err)
		}
	}

	f.Fuzz(ft)

}

func FuzzPatterns(f *testing.F) {

	pats := []string{
		"abc",
		"a|b",
		"a*d+",
		"^kjh$",
	}

	for _, pat := range pats {
		f.Add(pat)
	}

	ft := func(t *testing.T, pat string) {
		gen, err := NewGen(pat)
		if err != nil {
			return // skip invalid patterns that are not recognized
		}
		it := NewRandChooser()
		s := gen.Next(it)
		err = gen.Verify(s)
		if err != nil {
			t.Errorf("pattern : %#q, gen : %s, generated : %#q, error : %v", pat, gen, s, err)
		}
	}

	f.Fuzz(ft)
}

func TestCleanVisual(t *testing.T) {
	pats := []string{
		"0+?",
		"*$*",
	}

	for _, p := range pats {
		fmt.Printf("%#q\t-->\t%#q\n", p, cleanPattern(p))
	}
}
