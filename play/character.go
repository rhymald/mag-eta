package play

import (
	"rhymald/mag-eta/balance/primitives"
	"rhymald/mag-eta/balance/functions"
	"sync"
)

type Character struct {
	Base *BasicStats
	Atts *Attributes
	HP *primitives.Health
	sync.Mutex
}

type Attributes struct {
	ID int
	Vitality float64
}

type BasicStats struct {
	ID int
	Body *primitives.Stream
}

func Init_BasicStats() *BasicStats {
	var buffer BasicStats
	buffer.ID = functions.Epoch()
	buffer.Body = primitives.Init_Stream(primitives.PhysList[1])
	return &buffer
}

func (base *BasicStats) Init_Character() *Character {
	var buffer Character 
	buffer.Base = base
	buffer.HP = primitives.Init_Health()
	return &buffer
}

func (char *Character) Init_Attributes() {
	var buffer Attributes
	buffer.ID = functions.Epoch()
	char.Lock()
	buffer.Vitality = 10 + (*char.Base).Body.Mean()
	(*char).Atts = &buffer
	char.Unlock()
}