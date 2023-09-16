package world 

import (
	"sync"
	"math"
	"rhymald/mag-eta/balance/functions"
	"fmt"
)

type Grid struct {
	Sum [3]int // x, y, c
	Reg [][3]interface{}
	sync.Mutex
}

func Init_Grid(id string) *Grid { 
	buffer := &Grid{}
	init := make(map[string][2]int)
	init["0000-000000000-0-0000000"] = [2]int{0, 0}
	init[id] = [2]int{0, 0}
	buffer.Nonce(init)
	return buffer
}

func (gr *Grid) Nonce(write map[string][2]int) {
	gr.Lock()
	read := gr.GetAll(false)
	for id, pos := range write {
		if _, ok := read[id] ; ok {
			(*gr).Reg[read[id][2]][1] = (*gr).Reg[read[id][2]][1].(int) + pos[0]
			(*gr).Reg[read[id][2]][2] = (*gr).Reg[read[id][2]][2].(int) + pos[1]
		} else {
			buffer := [3]interface{}{ id, pos[0], pos[1] }
			(*gr).Reg = append((*gr).Reg, buffer)
		}
	}
	gr.Unlock()
}

func (gr *Grid) GetAll(lock bool) (map[string][3]int) { // x, y, i
	buffer := make(map[string][3]int)
	if lock { gr.Lock() }
	for i, each := range (*gr).Reg { if _, ok := buffer[each[0].(string)] ; !ok {
		buffer[each[0].(string)] = [3]int{ each[1].(int), each[2].(int), i } 
	} else {
		fmt.Println("WARNING[World.Grid.GetAll()] Duplicated id:", each[0].(string))
	}}
	if lock { gr.Unlock() }
	return buffer
}

func (gr *Grid) GetAgainst(step float64) [][]string { // step+ for all, delim by area; step- for within by +1
	if step <= 0.1 { step = 1 }
	read := gr.GetAll(true)
	var buffer = [][]string{ []string{} } 
	// var targetPos [2]int
	// ints, isInts := target.([2]int) 
	// strng, isStr := target.(string) ; isStr = isStr && len(strng) == 24
	for id, xyi := range read {
		far := functions.Round(math.Log2( 1 + math.Sqrt(float64(xyi[0]*xyi[0] + xyi[1]*xyi[1])) / 1000 / step))
		diff := far - len(buffer) + 1 
		if diff > 0 { for x:=0 ; x<diff ; x++ { buffer = append(buffer, []string{}) }}
		buffer[far] = append(buffer[far], id)
	}
	return buffer
}
