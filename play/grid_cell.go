package play

import (
	"sync"
	// "rhymald/mag-eta/balance/functions"
	// "fmt"
)


type Cell struct {
	TID map[int]*Registry
	sync.Mutex
}

func Init_Cell() *Cell { return &Cell{ TID: make(map[int]*Registry) }}

func (c *Cell) PUT(t int, id string) {
	// if c == nil { c = Init_Cell() }
	c.Lock()
	// buffer := (*c).TID[t]
	var buffer *Registry
	if _, ok := (*c).TID[t] ; !ok { (*c).TID[t] = Init_Registry() }
	buffer = (*c).TID[t]
	c.Unlock()
	buffer.Register(id)
}

func (c *Cell) GET(t int) map[string]int {
	c.Lock()
	if _, ok := (*c).TID[t] ; !ok { c.Unlock() ; return make(map[string]int) }
	buffer := (*c).TID[t]
	c.Unlock()
	list := buffer.Read()
	return list
}

// func (c *Cell) DELETE(ids []string) {
// 	c.Lock()
// 	buffer := (*c).TID
// 	for _, reg := range buffer { for _, id := range ids {
// 		reg.Deregister(id)
// 	}}
// 	(*c).TID = buffer
// 	c.Unlock()
// }
