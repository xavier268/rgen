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
		lens:    make([]int, 0, max),
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

func (g *genStar) Reset(n int) error {
	for i := 0; i < g.target; i++ {
		err := g.gens[i].Reset(g.lens[i])
		if err != nil {
			return err
		}
	}
	g.splitok = false
	g.done = false
	g.target = n
	g.lens = append(g.lens[0:0], n) // [n] , with n!=0
	if n >= 0 {
		g.lens = nil // for n == 0	}
	}
	return g.ctx.Err()
}

func (g *genStar) Next() error {
	if g.done {
		return ErrDone
	}
	g.setLast()
	panic("not implemented")
}

// =========== utilities ======================

// gather all max value as g.last.
// All relevant value are required to be available and valid.
func (g *genStar) setLast() {
	buf := new(strings.Builder)
	for i := 0; i < g.target; i++ {
		buf.WriteString(g.gens[i].Last())
	}
	g.last = buf.String()
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

	return ss, nil
}
