package world

import(
	"rhymald/mag-eta/play/character"
	"sync"
)

type ByIDList struct {
	List [][2]interface{}
	sync.Mutex
}

func Init_ByIDList() *ByIDList { return &ByIDList{} }

func (list *ByIDList) Read(id string) (*character.State, bool) {
	list.Lock()
	for _, kv := range (*list).List {
		k := kv[0].(string)
		if k == id { list.Unlock() ; return kv[1].(*character.State), true }
	}
	list.Unlock()
	return nil, false
}

func (list *ByIDList) Add(id string, st *character.State) {
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

func (list *ByIDList) GetAll() map[string]*character.State {
	buffer := make(map[string]*character.State)
	list.Lock()
	for _, kv := range (*list).List {
		buffer[kv[0].(string)] = kv[1].(*character.State)
	}
	list.Unlock()
	return buffer
}
