package play

import (
	"rhymald/mag-eta/balance/primitives"
	"rhymald/mag-eta/balance/functions"
	"sync"
)

type Character struct {
	Base *BasicStats
	Atts *Attributes
	HP   *primitives.Health
	Pool *primitives.Pool
	sync.Mutex
}

type Attributes struct {
	ID int
	Vitality float64
	Capacity float64
}

type BasicStats struct {
	ID int
	Body *primitives.Stream
	Energy []*primitives.Stream
}

func Init_BasicStats() *BasicStats {
	var buffer BasicStats
	buffer.ID = functions.Epoch()
	luck := functions.EpochNS() % 2
	for i:=0 ; i<3+luck ; i++ {
		buffer.Energy = append(
			buffer.Energy, 
			primitives.Init_Stream(primitives.ElemList[ (functions.EpochNS()%2)*(luck+1) ]),
		)
	}
	buffer.Body = primitives.Init_Stream(primitives.PhysList[1])
	return &buffer
}

func (base *BasicStats) Init_Character() *Character {
	var buffer Character 
	buffer.Base = base
	buffer.HP = primitives.Init_Health()
	buffer.Pool = primitives.Init_Pool()
	return &buffer
}

func (char *Character) Init_Attributes() {
	var buffer Attributes
	buffer.ID = functions.Epoch()
	char.Lock()
	buffer.Vitality = 10 
	buffer.Capacity = 15 
	(*char).Atts = &buffer
	char.Unlock()
}