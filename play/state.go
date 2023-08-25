package play

import (
	"rhymald/mag-eta/balance/primitives"
	// "rhymald/mag-eta/balance/functions"
	"sync"
)

type State struct {
	Trace [3]*Tracing
	Effects map[int]*primitives.Effect
	Write struct {
		Body primitives.Stream
		HP primitives.Health
		// Actions 
	}
	Current *Character
	Later struct {
		Body primitives.Stream
		HP primitives.Health
	}
	sync.Mutex
}

func (c *Character) Init_State() *State {
	var buffer State // race-4 race-7
	c.Lock()
	buffer.Current = c
	buffer.Effects = make(map[int]*primitives.Effect)
	// buffer.Later.Time = make(map[string]int)
	// buffer.Later.Time["Life"] = base.Epoch()
	buffer.Later.HP = *c.HP
	buffer.Later.Body = *c.Base.Body
	c.Unlock()
	// buffer.Writing.Time = make(map[string]int)
	// buffer.Writing.Time["Life"] = 0 
	// buffer.Writing.Life = *(base.MakeLife())
	// buffer.Writing.Life.Rate = 0
	for bucket := 0 ; bucket < 3 ; bucket++ {
		trace := Init_Tracing() 
		trace.Wipe()
		buffer.Trace[bucket] = trace
	}
	return &buffer
}
