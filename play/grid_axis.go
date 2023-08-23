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
	if _, ok := (*a).Idx[pos] ; !ok { (*a).Idx[pos] = Init_Cell() }
	buffer = (*a).Idx[pos]
	// add cell if null
	a.Unlock() 
	buffer.PUT(t, id)
}

func (a *Axis) Get(start, within, when int) map[string][3]int {
	a.Lock()
	read := (*a).Idx
	a.Unlock()
	buffer := make(map[string][3]int) // id: t, a, ts
	for i := 0; i < within/2; i++ {
		pos, neg := start-i, start+1
		for moment := when-1 ; moment < when+2 ; moment++ {
			if _, ok := read[pos]; !ok { continue } else {
				cell := read[pos].GET(moment)
				for id, ts := range cell {
					checkTS := [3]int{}
					if _, ok := buffer[id] ; ok { checkTS = buffer[id] }
					if checkTS[2] > ts { continue } else {
						buffer[id] = [3]int{ moment, pos, ts }
					}
				}
			}
			if _, ok := read[neg]; !ok { continue } else {
				cell := read[neg].GET(moment)
				for id, ts := range cell {
					checkTS := [3]int{}
					if _, ok := buffer[id] ; ok { checkTS = buffer[id] }
					if checkTS[2] > ts { continue } else {
						buffer[id] = [3]int{ moment, neg, ts }
					}
				}
			}
		}
	}
	return buffer
}