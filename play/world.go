package play

import (
	"rhymald/mag-eta/balance/functions"
	"sync"
)

type World struct {
	ID string
	Queue struct {
		Chan chan map[string][][3]int
		Buffer []map[string][][3]int
		sync.Mutex
	}
	Grid [3]*Grid
	ByID *ByIDList
}

type ByIDList struct {
	List map[string]*State
	sync.Mutex
}

func Init_ByIDList() *ByIDList { return &ByIDList{ List: make(map[string]*State) } }

func Init_World() *World {
	var buffer World
	buffer.ByID = Init_ByIDList()
	buffer.Queue.Chan = make(chan map[string][][3]int)
	buffer.Queue.Buffer = []map[string][][3]int{}
	buffer.ID = functions.GetID( functions.StartEpoch/1000000, functions.StartEpoch%1000000 )
	for i:=0 ; i<3 ; i++ {buffer.Grid[i] = Init_Grid(0, 0)}
	go func(){ (&buffer).GridWriter_ByPush() }()
	return &buffer
}

func (w *World) WhichGrid() (*Grid, *Grid) {
	tAxisStep, tRange := functions.TAxisStep, functions.TRange
	epoch := functions.Epoch()
	even := (epoch/(tRange*tAxisStep))%3
	read, write := (*w).Grid[(even+2)%3], (*w).Grid[(even+3)%3] 
	x, y := write.Get_CentralPos()
	(*w).Grid[(even+1)%3] = Init_Grid( x, y )
	return read, write
}

func (w *World) Login(st *State) string {
	w.ByID.Lock()
	id := functions.GetID((*st.Current.Base).ID, (*st.Current.Atts).ID)
	if _, ok := w.ByID.List[id] ; ok { w.ByID.Unlock() ; return "ERROR:AlreadyLoggedIn" } 
	w.ByID.List[id] = st
	w.ByID.Unlock()
	return id
}

func (w *World) GridWriter_ByPush() {
	writeToCache := (*w).Queue.Chan
	for {
		input := <- writeToCache // just a black hole
		for id, posList := range input { for _, pos := range posList {
			_, writer := w.WhichGrid()
			writer.Put_ID_to_XYT(id, pos[1], pos[2], pos[0])
		}}
	}
}
