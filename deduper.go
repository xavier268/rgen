package revregex

import (
	"crypto/md5"
	"math/big"
)

// Deduper is an object that can be used to deduplicate strings.
type Deduper interface {
	// Unique is true if string is seen for the first time.
	Unique(st string) bool
}

// -------------------------------------

// Creates a map-based Deduper.
// Memory footprint will grow indefinitely, but respose is always exact.
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
// True result are always right, but  after a while, False result are only most likely right.
// The larger the bloomsize, the less false positive.
func NewDedupBloom(bloomsize int) Deduper {
	d := new(dedupbloom)
	d.z = big.NewInt(0)
	d.bs = bloomsize
	return d
}

type dedupbloom struct {
	bs int      // bloomsize in bits
	z  *big.Int // we use this as a 256 bit field
}

// Reasonable default value up to a few hundred thousands different strings.
const DefaultBloomSize = 256 * 256 * 256 * 16 // total filter size in bits

func (d *dedupbloom) Unique(st string) bool {

	hh := md5.Sum([]byte(st)) // hh is 16 bytes long, or 16*256 bits long.
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
