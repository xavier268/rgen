package generator

// match nothing
type genNoMatch struct{}

var _ Generator = new(genNoMatch)

// Next implements Generator.
func (g *genNoMatch) Next() error {
	return ErrDone
}

// Reset implements Generator.
func (g *genNoMatch) Reset(length int) error {
	return nil
}

func (g *genNoMatch) Last() string {
	panic("last should never be called on genNoMatch")
}
