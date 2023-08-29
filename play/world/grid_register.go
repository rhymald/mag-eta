package world

import (
	"sync"
	"rhymald/mag-eta/balance/functions"
)

type Registry struct {
	IDT map[string]int 
	sync.Mutex
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

func (reg *Registry) Read() map[string]int {
	now := functions.Epoch()
	reg.Lock()
	read := (*reg).IDT
	reg.Unlock() 
	buffer := make(map[string]int)
	for id, old := range read {
		buffer[id] = old - now
	}
	return buffer
}
