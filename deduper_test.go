package revregex

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDedupMap(t *testing.T) {

	d := NewDedupMap()

	for i := 0; i < 10_000; i++ {
		if !d.Unique(fmt.Sprint(i)) {
			t.Fatal(i, "expecetd to be unique")
		}
	}

	for i := 0; i < 10_000; i++ {
		if d.Unique(fmt.Sprint(i)) {
			t.Fatal(i, "expected NOT to be unique")
		}
	}
}

func TestDedupBloom(t *testing.T) {

	d := NewDedupBloom(DefaultBloomSize)
	salt := "kjhs  sklhqskdhmlqsjfl!k"

	for i := 0; i < 5_000_000; i++ {
		if !d.Unique(fmt.Sprint(i, salt)) {
			t.Fatal(i, "expected to be unique - false positive !")
		}
	}

	for i := 0; i < 5_000_000; i++ {
		if d.Unique(fmt.Sprint(i, salt)) {
			t.Fatal(i, "expected NOT to be unique")
		}
	}
}

func TestDedupConstistentency(t *testing.T) {

	d1 := NewDedupBloom(DefaultBloomSize)
	d2 := NewDedupMap()
	rd := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 1; i < 1_000_000; i++ {

		k := rd.Int()
		ks := fmt.Sprint(k)
		if d1.Unique(ks) != d2.Unique(ks) {
			t.Fatalf("inconsistent after %d rounds ", i)
		}

	}

}
