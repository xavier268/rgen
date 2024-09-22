package dedup

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestUnique(t *testing.T) {

	// create map and bloom dedup
	b1, b2 := NewDedupMap(), NewDedupBloom(DefaultBloomSize)

	// generate random numbers, and ensure uniqueness response are the same.
	rd := rand.New(rand.NewSource(42))
	for i := 0; i < 10_000; i++ {
		tt := fmt.Sprint(rd.Intn(1_000))
		u1, u2 := b1.Unique(tt), b2.Unique(tt)
		if u1 != u2 {
			t.Error("Unique response should be the same")
		}
	}
	bb := b2.(*dedupbloom)
	fmt.Println("Actual bloom size in Bytes:", len(bb.z.Bytes()))
}
