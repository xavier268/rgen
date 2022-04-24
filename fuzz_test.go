package revregex

import "testing"

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
		g := NewGen(pat)
		it := NewBytesChooser(bb)
		r1 := g.Next(it)
		err := g.Verify(r1)
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
