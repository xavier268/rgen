package revregex

import (
	"math/big"
	"math/rand"
	"time"
)

// Chooser defines an interface that allows to make choices.
type Chooser interface {
	// Intn provides a number between 0 and (n-1) included.
	Intn(n int) int
}

// exp provides a number between 0 and infinity, where probability of a given length decreases exponentially.
// n=0 has a  probability 0.5, n=1 has a 0.25 probability, n=2 has a 0.125 prbaility, ...
func exp(it Chooser) (n int) {

	for {
		if it.Intn(2) == 0 {
			return n
		}
		n++
	}
}

// NewRandChooser uses random as the source for decision.
// It is garanteed that no string has a zero probability, but longuer strings have a much ower chance of appearing.
func NewRandChooser() Chooser {
	return NewRandChooserSeed(time.Hour.Milliseconds())
}

// NewRandChooserSeed uses random as the source for decision.
// It is garanteed that no string has a zero probability, but longuer strings have a much ower chance of appearing.
// Setting the seed allows for reproductibility in tests.
func NewRandChooserSeed(seed int64) Chooser {
	return rand.New(rand.NewSource(seed))
}

// NewBytesChooser uses buf as a source for decision. This makes the exploration of all possible strings perfectily deterministic.
func NewBytesChooser(buf []byte) Chooser {
	bb := new(bigChooser)
	bb.big = big.NewInt(0)
	bb.big.SetBytes(buf)

	return bb
}

type bigChooser struct {
	// big is always positive (or zero).
	big *big.Int
}

func (b *bigChooser) Intn(n int) int {
	if b == nil || n == 0 {
		return 0
	}
	r := new(big.Int)
	b.big.DivMod(b.big, big.NewInt(int64(n)), r)
	return int(r.Int64())
}
