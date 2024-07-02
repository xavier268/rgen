package rgen

import (
	"context"
	"regexp/syntax"
	"strings"
)

type genStar struct {
	*generator
	lens    []int // length expected from each generators
	target  int   // target length to generate - 0 alone is a valid split, or a split without any zero.
	splitok bool  // has the current split already been initialized with valid values ?
}

func newGenStar(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {
	g := &genStar{
		generator: &generator{
			ctx:  ctx,
			max:  max,
			gens: make([]Generator, max),
			last: "",
			done: false,
		},
		lens:    nil,
		target:  0,
		splitok: false,
	}

	// construct 'max' generators.
	for i := 0; i < max; i++ {
		var err error
		g.gens[i], err = newGenerator(g.ctx, re.Sub0[0], g.max)
		if err != nil {
			return nil, err
		}
	}

	return g, ctx.Err()
}

// Just bare minimum. No strings initialized.
// lens starts with [n] or [].
// Only corresponding generator is reset.
func (g *genStar) Reset(n int) error {
	if n < 0 {
		panic("illegal negative length at reset")
	}
	// reset lens to [n]
	if n > 0 {
		g.lens = []int{n}
		g.gens[0].Reset(n) // reset first generator
	} else {
		g.lens = nil
	}
	g.splitok = false
	g.done = false
	g.target = n
	g.last = ""
	return g.ctx.Err()
}

func (g *genStar) Next() error {
	if g.done {
		return ErrDone
	}
	if g.target == 0 { // return "" only once
		g.done = true
		g.last = ""
		return g.ctx.Err()
	}

	if !g.splitok {
		if err := g.useNewSplit(); err != nil {
			return err
		}
		g.splitok = true
		g.setLast()
		// fmt.Println("DEBUG split:", g.lens)
		return g.ctx.Err()
	}

	// now, g.splitok = true
	// here we already have an initialized split, with valid values.
	// lets try to increment ?
	err := g.doNext(0)
	if err != nil {
		// we could not increment, we need another split.
		err := g.useNewSplit()
		if err != nil {
			return err
		}
		// ok - we have a new useable split.
		g.setLast()
		// fmt.Println("DEBUG split:", g.lens)
		return g.ctx.Err()
	}
	// here, we could increment within the existing split.
	// fmt.Println("DEBUG split:", g.lens)
	g.setLast()
	return g.ctx.Err()
}

// =========== utilities ======================

// try to execute next, or return error
// this will never change the split.
// It will only look at gens[i:].
// g.last is not updated.
// it assumes all first values of the split have been initialized.
func (g *genStar) doNext(i int) error {

	if i >= len(g.lens) {
		return ErrDone
	}

	// try to use i
	if g.gens[i].Next() == nil {
		// success, return
		return nil
	} else {
		// reset i
		err := g.gens[i].Reset(g.lens[i])
		if err != nil {
			return err
		}
		// retrieve first value
		err = g.gens[i].Next()
		if err != nil {
			return err
		}
	}

	// since we could not use i, try to increment i+1
	return g.doNext(i + 1)
}

// gather all max value as g.last.
// All relevant value are required to be available and valid.
func (g *genStar) setLast() {
	buf := new(strings.Builder)
	for i := 0; i < len(g.lens); i++ {
		buf.WriteString(g.gens[i].Last())
	}
	g.last = buf.String()
}

// loop until a new valild split can be found.
// reset is then performed with this new split, and first values are initilized.
// error means no further split will ever work.
// if slitok is false, the current split is not changed, but just reset and initialized.
func (g *genStar) useNewSplit() (err error) {
	if g.splitok { // if current split was initialized, look for another one.
		g.lens, err = incSplitStar(g.lens)
		if err != nil {
			// no more split to try
			return ErrDone
		}
	}

	// here, we do not increment the split, but first try to use it.
	for i := range g.lens {
		gen := g.gens[i]
		if err := gen.Reset(g.lens[i]); err != nil {
			g.splitok = true       // force incrementing split
			return g.useNewSplit() // cannot reset this generator - try another split
		}
		if err := gen.Next(); err != nil {
			g.splitok = true       // force incrementing split
			return g.useNewSplit() // cannot initialize first value try another split
		}
	}
	// all good
	g.splitok = true
	return g.ctx.Err()
}

// Genere tous les splits dont la somme est constante.
// Le split initial est [n] (n != 0)
// ou [] pour une somme nulle.
// chaque terme du split doit être STRICTEMENT positif.
// Il en resulte que le split ne peut dépasser la longueur n.
// Si la slice a une capacité de max >= n, les changements
// in place ne posent pas de problème.
func incSplitStar(s []int) ([]int, error) {

	if len(s) == 0 {
		return nil, ErrDone
	}

	// on cherche le dernier nombre > 1
	ttl := 0
	last := -1
	for i, v := range s {
		ttl += v // calcul de la somme
		if v > 1 {
			last = i // capture de last
		}
	}
	if last == -1 {
		return nil, ErrDone
	}

	// make a copy of s
	res := make([]int, len(s))
	copy(res, s)

	// On decrement le [last], et on regroupe tout ce qui suit dans une seul nombre, pour garder le total inchangé.
	res[last] = s[last] - 1
	ss := append(res[:last+1], 1+sum(res[last+1:]))

	//fmt.Println("DEBUG split : ", ss)
	return ss, nil
}
