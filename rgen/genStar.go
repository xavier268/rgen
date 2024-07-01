package rgen

import (
	"errors"
	"fmt"
)

type genStar struct {
	*generator
	lens    []int // length expected from each generators
	target  int   // target length to generate - 0 alone is a valid split, or a split without any zero.
	splitok bool  // has the current split already been initialized with valid values ?

}

// =========== utilities ======================

// Genere tous les splits dont la somme est constante.
// Le split initial est [n] (n != 0)
// ou [] pour une somme nulle.
// chaque terme du split doit être STRICTEMENT positif.
// Il en resulte que le split ne peut dépasser la longueur n.
// Si la slice a une capacité de max >= n, les changements
// in place ne posent pas de problème.
func incSplitStar(s []int) ([]int, error) {

	if len(s) == 0 {
		return nil, errors.New("toutes les combinaisons ont été générées")
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
	fmt.Println("last : ", last)
	if last == -1 {
		return nil, errors.New("toutes les combinaisons ont été générées")
	}
	// make a copy of s
	res := make([]int, len(s))
	copy(res, s)

	// On decrement le [last], et on regroupe tout ce qui suit dans une seul nombre, pour ajuster le total.
	res[last] = s[last] - 1

	fmt.Println("keep :", res[:last+1])
	ss := append(res[:last+1], 1+sum(res[last+1:]))

	return ss, nil
}
