package play

import (
	"sync"
	"rhymald/mag-eta/balance/functions"
	"math"
)

type Grid struct {
	Center struct {
		X int
		Y int
	}
	W struct {
		Counter int 
		Sum [2]int
		sync.Mutex
	}
	X, Y, V, R *Axis
	Reg *Registry
}

func Init_Grid(x, y int) *Grid {
	var buffer Grid
	buffer.X = Init_Axis()
	buffer.Y = Init_Axis()
	buffer.R = Init_Axis()
	buffer.V = Init_Axis()
	buffer.Reg = Init_Registry()
	buffer.W.Sum = [2]int{ 0, 0 }
	buffer.W.Counter = 0
	buffer.Center.X, buffer.Center.Y = x, y 
	return &buffer
}
func (gr *Grid) Get_CentralPos() (int, int) {
	x, y := 0, 0
	(*gr).W.Lock()
	if (*gr).W.Counter == 0 { (*gr).W.Unlock() ; return 0, 0 }
	x, y = (*gr).W.Sum[0] / (*gr).W.Counter, (*gr).W.Sum[1] / (*gr).W.Counter
	(*gr).W.Unlock()
	return x, y
}

func (gr *Grid) Put_ID_to_XYT(id string, x, y, t int) {
	xc, yc := gr.Get_CentralPos()
	x += -xc ; y += -yc
	r := functions.Round( math.Sqrt( float64(x*x + y*y) ))
	v := functions.Round( math.Atan( float64(y)/float64(x) ) / math.Pi * 1000 ) 
	gr.X.Put(t, x, id)
	gr.Y.Put(t, x, id)
	gr.V.Put(t, v, id)
	gr.R.Put(t, r, id)
	gr.W.Lock()
	(*gr).W.Counter += 1
	write := (*gr).W.Sum
	write[0] += x ; write[1] += y
	(*gr).W.Sum = write
	gr.W.Unlock()
	gr.Reg.Register(id)
}

func (gr *Grid) Get_Square(x, y, r int) map[string][4]int {
	buffer := make(map[string][4]int) // id: x, y, t, ts
	return buffer
}
// + Get_Round
// + Get_Sector
