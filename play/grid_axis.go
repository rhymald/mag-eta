package play

import (
	"sync"
	// "rhymald/mag-eta/balance/functions"
)

type Axis struct {
	Idx map[int]*Cell 
	sync.Mutex
}

func Init_Axis() *Axis { return &Axis{ Idx: make(map[int]*Cell) }}

func (a *Axis) Put(t, pos int, id string) {
	a.Lock()
	var buffer *Cell
	if _, ok := (*a).Idx[pos] ; !ok { (*a).Idx[pos] = Init_Cell() } // race-1
	buffer = (*a).Idx[pos]
	a.Unlock() 
	buffer.PUT(t, id)
}

func (a *Axis) Get(start, within, when int) map[string][3]int {
	a.Lock()
	buffer := make(map[string][3]int) // id: t, a, ts
	read := (*a).Idx
	a.Unlock()
ByPositionOnAxis:
	for pos:=start-within/2 ; pos<=start+within/2 ; pos++ {
		cell, ok := read[pos] ; if !ok { continue ByPositionOnAxis }
	ByTimeAxis:
		for t:=when-4 ; t<=when+1 ; t++ {
			cellTcontent := cell.GET(t)
			if len(cellTcontent) == 0 { continue ByTimeAxis }
		ByWriteTimeStamp:
			for id, ts := range cellTcontent {
				compare, ook := buffer[id]
				if !ook {
					buffer[id] = [3]int{ t, pos, ts }
				} else {
					if compare[2] < ts { buffer[id] = [3]int{ t, pos, ts } } else { continue ByWriteTimeStamp }
				}
			} // end ByWriteTimeStamp
		} // end ByTimeAxis
	} // end ByPositionOnAxis
	return buffer
}
