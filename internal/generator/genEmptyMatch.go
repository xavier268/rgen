package generator

// only match ""
type genEmptyMatch struct {
	done bool
}

var _ Generator = new(genEmptyMatch)

// Next implements Generator.
func (g *genEmptyMatch) Next() error {
	if g.done {
		return ErrDone
	}
	g.done = true
	return nil
}

// Reset implements Generator.
func (g *genEmptyMatch) Reset(length int) error {
	g.done = (length != 0)
	return nil
}

func (g *genEmptyMatch) Last() string {
	return ""
}
