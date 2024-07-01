package rgen

/*
type genConcat struct {
	*generator
	frag0      *string // current first fragment, nil if not set
	len0, len1 int     // both length, initially (n,0)
}

func newGenConcat(ctx context.Context, re *syntax.Regexp, max int) (Generator, error) {
	if len(re.Sub) != 2 {
		panic("concat expects exactly 2 arguments once preprocessed")
	}
	gen1, err := newGenerator(ctx, re.Sub[0], max)
	if err != nil {
		return nil, err
	}
	gen2, err := newGenerator(ctx, re.Sub[1], max)
	if err != nil {
		return nil, err
	}

	return &genConcat{
		generator: &generator{
			ctx:  ctx,
			max:  max,
			gens: []Generator{gen1, gen2},
			done: false,
		},
		frag0: nil,
		len0:  0,
		len1:  0,
	}, ctx.Err()

}

func (g *genConcat) Reset(n int) error {
	err := g.generator.Reset(n)
	if err != nil {
		return err
	}
	g.len0, g.len1 = n, 0
	g.frag0 = nil
	return g.ctx.Err()
}

func (g *genConcat) Next() (string, error) {

	if g.done {
		return "", ErrDone
	}

	if g.len0 < 0 {
		g.done = true
		return "", ErrDone
	}

	// TODO $$$$$$$$$$$$$$$$$$

	// if frag0 is nil, try to set it and loop

	// get next frag from gen2. If we could, fine, send result.

	// If we cant, change split, reset generators, reset frag0 to nil, and loop

	panic("todo")
}

// ==================================================================

// replace abcd by a(b(c(d ...))) with always exctly 2 arguments to OpConcat.
// called by newGenerator(...) for every regexp, before further processing.
func processConcat(re *syntax.Regexp) *syntax.Regexp {

	if re.Op != syntax.OpConcat {
		return re
	}

	if len(re.Sub) == 0 {
		panic("concat with 0 arguments")
	}
	if len(re.Sub) == 1 {
		return re.Sub[0]
	}
	if len(re.Sub) == 2 {
		return re // no change
	}

	// here, more than 2 args
	return &syntax.Regexp{
		Op: syntax.OpConcat,
		Sub: []*syntax.Regexp{
			re.Sub[0],
			processConcat(&syntax.Regexp{
				Op:  syntax.OpConcat,
				Sub: re.Sub[1:],
			}),
		},
	}
}

*/
