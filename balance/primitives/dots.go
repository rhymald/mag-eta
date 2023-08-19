package primitives

type Dot struct {
	E string
	W int
}

func (dot *Dot) Weight() float64 { return float64((*dot).W)/1000 }