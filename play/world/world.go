package world

import (
	"rhymald/mag-eta/balance/functions"
	"rhymald/mag-eta/play/character"
	// "math"
	"sync"
	"fmt"
)

type World struct {
	ID string
	Grid *Grid
	Queue struct {
		Chan chan map[string][2]int
		Buffer []map[string][2]int
		sync.Mutex
		Size int
		Timeout int
	}
	ByID *ByIDList
	sync.Mutex
}

func Init_World() *World {
	var buffer World
	buffer.ByID = Init_ByIDList()
	buffer.Queue.Chan = make(chan map[string][2]int)
	buffer.Queue.Buffer = []map[string][2]int{}
	buffer.Queue.Size = 2 * functions.TRange * 1000 / 618
	buffer.Queue.Timeout = functions.TAxisStep * 1000 / 618
	buffer.ID = functions.GetID( functions.StartEpoch/3600000000, functions.StartEpoch%3600000000 )
	buffer.Grid = Init_Grid(buffer.ID)
	// for i:=0 ; i<3 ; i++ {buffer.Grid[i] = Init_Grid(0, 0)}
	// go func(){ (&buffer).GridWriter_FromBuffer() }()
	go func(){ (&buffer).GridBuffer_ByPush() }()
	return &buffer
}

// func (w *World) WhichGrid() (*Grid, *Grid) {
// 	tAxisStep, tRange := functions.TAxisStep, functions.TRange
// 	epoch := functions.Epoch()
// 	even := (epoch/(tRange*tAxisStep))%3
// 	w.Lock()
// 	read, write := (*w).Grid[(even+2)%3], (*w).Grid[(even+3)%3] 
// 	x, y := write.Get_CentralPos()
// 	// (*w).Grid[(even+1)%3] = Init_Grid( x, y )
// 	w.Unlock()
// 	return read, write
// }

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

func (w *World) GridBuffer_ByPush() {
	// var wg sync.WaitGroup
	_, timewatcher := 0, 0
	timer := 0 
	writeToCache := (*w).Queue.Chan
	for { //for input := range writeToCache {
		input := <- writeToCache // just a black hole
		(*w).Grid.GetAgainst(0)
		(*w).Queue.Lock()
		(*w).Queue.Buffer = append((*w).Queue.Buffer, input)
		triggered := functions.Epoch() - timer >= (*w).Queue.Timeout
		bufferSize := len((*w).Queue.Buffer) 
		triggered = triggered || bufferSize >= (*w).Queue.Size
		(*w).Queue.Unlock()
		if triggered {
			timer = functions.Epoch()
			(*w).Queue.Lock()
			buffer := (*w).Queue.Buffer
			(*w).Queue.Buffer = []map[string][2]int{}
			(*w).Queue.Unlock()
			write := make(map[string][2]int)
			for _, each := range buffer { for id, pos := range each {
				write[id] = pos
			}}
			(*w).Grid.Nonce(write)
			timewatcher = functions.EpochNS() - timer*1000000
			fmt.Printf("\r => Written: %9.3fms / %4d = %9.3fms\r", float64(timewatcher)/1000000, bufferSize, float64(timewatcher)/1000000/float64(bufferSize))
		}
		// for id, posList := range input { 
		// 	wg.Add(1)
		// 	start := functions.EpochNS()
		// 	for _, pos := range posList {
		// 		_, writer := w.WhichGrid()
		// 		writer.Put_ID_to_XYT(id, pos[1], pos[2], pos[0])
		// 	}
		// 	list := w.Seek_Square( posList[len(posList)-1][1], posList[len(posList)-1][2], 1400 )
		// 	if len(list) > 1 { 
		// 		fmt.Printf(" ==> RW time: %0.3fms\n", -(float64(start)+float64(functions.EpochNS()))/1000000 )
		// 		fmt.Println("    ", list) 
		// 	}
		// 	wg.Done()
		// }
		// wg.Wait()
	}
}
// func (w *World) GridWriter_FromBuffer() {
// 	pause := 0.0// float64(functions.TAxisStep) / math.Phi
// 	var wg sync.WaitGroup
// 	for {
// 		(*w).Queue.Lock()
// 		input := (*w).Queue.Buffer
// 		if len(input) < 10 { (*w).Queue.Unlock() ; continue }
// 		(*w).Queue.Buffer = []map[string][][3]int{}
// 		(*w).Queue.Unlock()
// 		wg.Add(1)
// 		go func(wg *sync.WaitGroup){
// 			counterT, found, avgT, meanT := 0, 0, 0.0, 0.0
// 			for _, each := range input { for id, posList := range each {
// 				start := functions.EpochNS()
// 				for _, pos := range posList {
// 					_, writer := w.WhichGrid()
// 					writer.Put_ID_to_XYT(id, pos[1], pos[2], pos[0])
// 				}
// 				avgT += float64(functions.EpochNS())/1000000-float64(start)/1000000
// 				counterT++
// 				meanT += 1000000/float64(functions.EpochNS()-start)
// 				list := w.Seek_Square( posList[len(posList)-1][1], posList[len(posList)-1][2], 1400 )
// 				found = len(list)
// 			}}
// 			if counterT != 0 {fmt.Printf("\r                 ==> RW time[%d/%d]:\tmean=%0.3fms\tavg=%0.3fms\ttotal=%0.3fms \r", found, counterT, float64(counterT)/meanT, avgT/float64(counterT), avgT )}
// 			wg.Done()
// 		}(&wg)
// 		functions.Wait( pause )
// 	}
// 	wg.Wait()
// }

// func (w *World) Seek_Square(x, y, r int) []string {
// 	reader, writer := w.WhichGrid()
// 	buffer := writer.Get_Square(x, y, r)
// 	old := reader.Get_Square(x, y, r)
// 	for id, row := range old { if _, ok := buffer[id] ; !ok { 
// 		row[2] += -functions.TRange
// 		buffer[id] = row
// 	}}
// 	list := []string{}
// 	states := (*w).ByID.GetAll()
// 	for t:=functions.TRange-1 ; t>=-functions.TRange ; t-- { for id, row := range buffer {
// 		player, ok := states[id]
// 		path := player.Path()
// 		actual := row[0] == path[1][0] && row[1] == path[1][1] && ok
// 		if row[2] == t && actual {
// 			list = append(list, fmt.Sprintf("id = %s, x = %6d, y = %6d, t = %3d, old = %6dms", id, row[0], row[1], row[2], row[3]))
// 		}
// 	}}
// 	return list
// }
