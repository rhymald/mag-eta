package primitives

import (
	"rhymald/mag-eta/balance/functions"
)

// main struct
type Effect struct {
	Time int `json:"Time"`
	Collision [2]int `json:"Collision"`
	Effects []Different_Effects `json:"Effects"`
}
type Different_Effects interface{
	Delayed() int
	HP() int
	Dots() (int, []Dot)
}


func Init_Effect() *Effect {
	buffer := Effect{ Time: functions.Epoch() }
	return &buffer
}


// delayed ============================================================
type Effect_MakeDot struct {
	Dot Dot `json:"Dot"`
	Delay int `json:"Delay"`
}
func (md Effect_MakeDot) Delayed() int { return md.Delay }
func (md Effect_MakeDot) HP() int { return 0 }
func (md Effect_MakeDot) Dots() (int, []Dot) { return md.Delay, []Dot{ md.Dot } }
func (ef *Effect) Add_Self_MakeDot(dot *Dot) float64 { 
	(*ef).Effects = append((*ef).Effects, 
	Effect_MakeDot{ 
		Dot: *dot,
		Delay: functions.CeilRound(1618/dot.Weight()+1),
	})
	return 1618/dot.Weight()+1
}


// instant ============================================================
type Effect_HPRegen struct {
	Portion int `json:"Portion"`
}
func (md Effect_HPRegen) Delayed() int { return 0 }
func (md Effect_HPRegen) HP() int { return md.Portion }
func (md Effect_HPRegen) Dots() (int, []Dot) { return 0, []Dot{} }
func (ef *Effect) Add_Self_HPRegen(hp int) { (*ef).Effects = append((*ef).Effects, 
	Effect_HPRegen{ 
		Portion: hp,
	})
}


// conditions ============================================================
// ...