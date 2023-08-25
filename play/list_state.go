package play

import(
	"sync"
)

type ByIDList struct {
	// List map[string]*State
	List [][2]interface{}
	sync.Mutex
}

func Init_ByIDList() *ByIDList { return &ByIDList{} }

func (list *ByIDList) Read(id string) (*State, bool) {
	list.Lock()
	for _, kv := range (*list).List {
		k := kv[0].(string)
		if k == id { list.Unlock() ; return kv[1].(*State), true }
	}
	list.Unlock()
	return nil, false
}

func (list *ByIDList) Add(id string, st *State) {
	list.Lock()
	(*list).List = append((*list).List, [2]interface{}{ id, st })
	list.Unlock()
}

func (list *ByIDList) Len() int {
	list.Lock()
	leng := len((*list).List)
	list.Unlock()
	return leng 
}

func (list *ByIDList) GetAll() map[string]*State {
	buffer := make(map[string]*State)
	list.Lock()
	for _, kv := range (*list).List {
		buffer[kv[0].(string)] = kv[1].(*State)
	}
	list.Unlock()
	return buffer
}
