package play

import (
	"rhymald/mag-eta/balance/primitives"
)

type State struct {
	// Effects
	Write struct {
		Body primitives.Stream // Body
		HP primitives.Health// Health 
		// Actions 
	}
	Current *Character
	Later Character
}