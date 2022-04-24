package revregex

import (
	"math/big"
	"math/rand"
	"time"
)

type Inter interface {
	// Intn provides a number between 0 and (n-1) included.
	Intn(n int) int
}

// exp provides a number between 0 and infinity, where probability of a given length decreases exponentially.
// n=0 has a  probability 0.5, n=1 has a 0.25 probability, n=2 has a 0.125 prbaility, ...
func exp(it Inter) (n int) {

	for {
		if it.Intn(2) == 0 {
			return n
		}
		n++
	}
}

func NewRandInter() Inter {
	return rand.New(rand.NewSource(time.Hour.Milliseconds()))
}

func NewBytesInter(buf []byte) Inter {
	bb := new(bigInter)
	bb.big = big.NewInt(0)
	bb.big.SetBytes(buf)

	return bb
}

type bigInter struct {
	// bif is always positive or zero.
	big *big.Int
}

func (b *bigInter) Intn(n int) int {
	if b == nil || n == 0 {
		return 0
	}
	r := new(big.Int)
	b.big.DivMod(b.big, big.NewInt(int64(n)), r)
	return int(r.Int64())
}
