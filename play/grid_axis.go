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
	if _, ok := (*a).Idx[pos] ; !ok { (*a).Idx[pos] = Init_Cell() } // here
	buffer = (*a).Idx[pos]
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
				cell := read[pos].GET(moment) // here
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

// WARNING: DATA RACE
// Write at 0x00c002c3ddd0 by goroutine 7:
//   runtime.mapassign_fast64()
//       /snap/go/10319/src/runtime/map_fast64.go:93 +0x0
//   rhymald/mag-eta/play.(*Axis).Put()
//       /home/eosipov/Documents/mag-eta/play/grid_axis.go:18 +0x13d
//   rhymald/mag-eta/play.(*Grid).Put_ID_to_XYT()
//       /home/eosipov/Documents/mag-eta/play/grid_grid.go:49 +0x33a
//   rhymald/mag-eta/play.(*World).GridWriter_ByPush()
//       /home/eosipov/Documents/mag-eta/play/world.go:65 +0x19a
//   rhymald/mag-eta/play.Init_World.func1()
//       /home/eosipov/Documents/mag-eta/play/world.go:34 +0x2e

// Previous read at 0x00c002c3ddd0 by goroutine 10:
//   runtime.mapaccess1_fast64()
//       /snap/go/10319/src/runtime/map_fast64.go:13 +0x0
//   rhymald/mag-eta/play.(*Axis).Get()
//       /home/eosipov/Documents/mag-eta/play/grid_axis.go:34 +0x1af
//   rhymald/mag-eta/api.testWorld()
//       /home/eosipov/Documents/mag-eta/api/testing.go:24 +0x16c
//   github.com/gin-gonic/gin.(*Context).Next()
//       /home/eosipov/go/pkg/mod/github.com/gin-gonic/gin@v1.9.1/context.go:174 +0x211
//   github.com/gin-gonic/gin.LoggerWithConfig.func1()
