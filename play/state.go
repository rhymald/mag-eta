package play

import (
	"rhymald/mag-eta/balance/primitives"
	"sync"
)

func CleanTrace() *Tracing { return &Tracing{ Trxy: make(map[int][3]int) }}
type Tracing struct {
	Trxy map[int][3]int
	sync.Mutex
}

type State struct {
	Trace [3]*Tracing
	// Effects
	Write struct {
		Body primitives.Stream
		HP primitives.Health
		// Actions 
	}
	Current *Character
	Later Character
}

