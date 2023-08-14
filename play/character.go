package play

import (
	"rhymald/mag-eta/balance/primitives"
	"rhymald/mag-eta/balance/balance"
)

type Character struct {
	HP int
	Atts *Attributes
	Base *BasicStats
}

type Attributes struct {
	Vitality float64
}

type BasicStats struct {
	Body *primitives.Stream
}

func InitStats() *BasicStats {
	var buffer BasicStats
	return &buffer
}