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
	}
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
	go func(){ (&buffer).GridWriter_ByPush() }()
	return &buffer
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
		_ = <- writeToCache
		// (*w).Queue.Buffer = append((*w).Queue.Buffer, char)
	}
}
