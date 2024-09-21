package rgen

import (
	"crypto/md5"
	"iter"
	"math/big"
)

// Deduper is an object that can be used to deduplicate strings.
type Deduper interface {
	// If Unique is true,  the string is garanteed to be seen for the first time.
	// If Unique is false, the string was most probably already seen.
	Unique(st string) bool
}

// Deduplicate the provided string iterator, using the provided Deduper.
// Returns a new string iterator, with all duplicated strings skipped.
// The order of the elements is guaranteed to match the initial iterator.
func Dedup(it iter.Seq[string], d Deduper) iter.Seq[string] {

	return func(yield func(string) bool) {
		for st := range it {
			if d.Unique(st) {
				if !yield(st) {
					return
				}
			} // else skip duplicates and iterate ...
		}
	}
}

// -------------------------------------

// Creates a map-based Deduper.
// Unique always provides an exact response.
// Memory footprint will grow infinitely, but response is always exact.
// Avoid in production.
func NewDedupMap() Deduper {
	d := new(dedupmap)
	d.m = make(map[string]bool, 10)
	return d
}

type dedupmap struct {
	m map[string]bool
}

func (d *dedupmap) Unique(st string) bool {
	if d.m[st] {
		return false
	}
	d.m[st] = true
	return true

}

// ----------------------------------------

// Creates a bloom-filter based Deduper.
// Memory footprint is fixed, driven by bloomsize setting in bits.
// True "Unique"" results are always right, but  after a while, False Unique could be wrong (appear as non unique, despite really being unique).
// The larger the bloomsize, the lesser the error rate.
func NewDedupBloom(bloomsize int) Deduper {
	d := new(dedupbloom)
	d.z = big.NewInt(0)
	d.bs = bloomsize
	return d
}

type dedupbloom struct {
	bs int      // bloomsize in bits
	z  *big.Int // we use this big.Int as a 256 bit field
}

// Reasonable default value up to a few hundred thousands different strings (appx 33M bytes memory foot print)
const DefaultBloomSize = 256 * 256 * 256 * 16 // total filter size

func (d *dedupbloom) Unique(st string) bool {

	hh := md5.Sum([]byte(st)) // hh is 16 bytes.
	// On average, we will set about 16 bits at the beginning ...
	changed := false
	bb := 0
	for i := 1; i < len(hh); i++ {
		bb = (bb*17 + 13*int(hh[i]) + 11) % d.bs // which bit to set ?
		if d.z.Bit(bb) == 0 {
			d.z.SetBit(d.z, bb, 1)
			changed = true
		}

	}
	return changed
}
