package primitives

import (
	"sync"
	"rhymald/mag-eta/balance/functions"
	"math"
)

type Dot struct {
	E string
	W int
}

type Pool struct {
	Dots map[int]*Dot
	sync.Mutex
}

func (dot *Dot) Weight() float64 { return float64((*dot).W)/1000 }
func Init_Pool() *Pool { return &Pool{ Dots: make(map[int]*Dot) }}


func (str *Stream) Init_Dot() *Dot { 
	a := math.Log2(str.Mean()+2)/math.Log2(7) 
	w := math.Pow(a,a) 
	return &Dot{ E: (*str).E, W: functions.FloorRound(functions.Ntrp( w*1000 )) }
}