package play

import (
	"sync"
	"rhymald/mag-eta/balance/functions"
)

type Cell struct {
	TID map[int][]string 
	sync.Mutex
}
type Axis struct {
	Idx map[int]*Cell 
	sync.Mutex
}
type Registry struct {
	IDT map[string]int 
	sync.Mutex
}

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

func Init_Registry() *Registry { return &Registry{ IDT: make(map[string]int) }}
func (reg *Registry) Register(id string) {
	reg.Lock()
	buffer := (*reg).IDT
	buffer[id] = functions.Epoch()
	(*reg).IDT = buffer
	reg.Unlock() 
}
func (reg *Registry) Deregister(id string) {
	reg.Lock()
	buffer := (*reg).IDT
	delete(buffer, id)
	(*reg).IDT = buffer
	reg.Unlock() 
}

func Init_Cell() *Cell { return &Cell{ TID: make(map[int][]string) }}

func Init_Axis() *Axis { return &Axis{ Idx: make(map[int]*Cell) }}

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
func (gr *Grid) CentralPos() (int, int) {
	x, y := 0, 0
	(*gr).W.Lock()
	if (*gr).W.Counter == 0 { (*gr).W.Unlock() ; return 0, 0 }
	x, y = (*gr).W.Sum[0] / (*gr).W.Counter, (*gr).W.Sum[1] / (*gr).W.Counter
	(*gr).W.Unlock()
	return x, y
}