package world

import (
	"rhymald/mag-eta/balance/functions"
	"rhymald/mag-eta/play/character"
	"sync"
	"fmt"
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
	sync.Mutex
}

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
	w.Lock()
	read, write := (*w).Grid[(even+2)%3], (*w).Grid[(even+3)%3] 
	x, y := write.Get_CentralPos()
	(*w).Grid[(even+1)%3] = Init_Grid( x, y )
	w.Unlock()
	return read, write
}

func (w *World) Login(st *character.State) string {
	// w.ByID.Lock()
	(*st.Current).Base.Lock() ; (*st.Current).Atts.Lock()
	id := functions.GetID((*st.Current.Base).ID, (*st.Current.Atts).ID)
	(*st.Current).Base.Unlock() ; (*st.Current).Atts.Unlock()
	if _, ok := (*w).ByID.Read(id) ; ok { return "ERROR:AlreadyLoggedIn" } 
	(*w).ByID.Add(id, st) // race-3 race-8
	// w.ByID.Unlock()
	return id
}

func (w *World) GridWriter_ByPush() {
	writeToCache := (*w).Queue.Chan
	for {
		input := <- writeToCache // just a black hole
		for id, posList := range input { 
			start := functions.EpochNS()
			for _, pos := range posList {
				_, writer := w.WhichGrid()
				writer.Put_ID_to_XYT(id, pos[1], pos[2], pos[0])
			}
			list := w.Seek_Square( posList[len(posList)-1][1], posList[len(posList)-1][2], 1400 )
			if len(list) > 1 { 
				fmt.Printf(" ==> RW time: %0.3fms\n", -(float64(start)+float64(functions.EpochNS()))/1000000 )
				fmt.Println("    ", list) 
			}
		}
	}
}

func (w *World) Seek_Square(x, y, r int) []string {
	reader, writer := w.WhichGrid()
	buffer := writer.Get_Square(x, y, r)
	old := reader.Get_Square(x, y, r)
	for id, row := range old { if _, ok := buffer[id] ; !ok { 
		row[2] += -functions.TRange
		buffer[id] = row
	}}
	list := []string{}
	states := (*w).ByID.GetAll()
	for t:=functions.TRange-1 ; t>=-functions.TRange ; t-- { for id, row := range buffer {
		player, ok := states[id]
		path := player.Path()
		actual := row[0] == path[1][0] && row[1] == path[1][1] && ok
		if row[2] == t && actual {
			list = append(list, fmt.Sprintf("id = %s, x = %6d, y = %6d, t = %3d, old = %6dms", id, row[0], row[1], row[2], row[3]))
		}
	}}
	return list
}