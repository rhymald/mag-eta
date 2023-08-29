package character 

import (
	"rhymald/mag-eta/balance/functions"
	"math"
	"sync"
)

type Tracing struct {
	Trxy [][4]int
	sync.Mutex
}

func Init_Tracing() *Tracing { return &Tracing{}}// Trxy: make(map[int][3]int) }} // race-5 race-7
func (tr *Tracing) Wipe() {  tr.Lock() ; (*tr).Trxy = [][4]int{} ; tr.Unlock() }

func (st *State) WhichTrace() (*Tracing, *Tracing, int) {
	tAxisStep, tRange := functions.TAxisStep, functions.TRange
	epoch := functions.Epoch()
	even := (epoch/(tRange*tAxisStep))%3
	st.Lock() // race-4
	read, write := (*st).Trace[(even+2)%3], (*st).Trace[(even+3)%3] 
	(*st).Trace[(even+1)%3] = Init_Tracing() // race-7
	wipe := (*st).Trace[(even+1)%3] 
	wipe.Wipe()
	st.Unlock()
	return read, write, epoch
}

func (st *State) Move(rotate float64, step bool, writeToCache chan map[string][][3]int) {
	tAxisStep, tRange := functions.TAxisStep, functions.TRange
	// epoch := functions.Epoch()
	// (*st).Trace[0].Lock() ; (*st).Trace[1].Lock() ; (*st).Trace[2].Lock()
	// later, trace, wipe := (*st).Trace[(even+2)%3].Trxy, (*st).Trace[(even+3)%3].Trxy, &(*st).Trace[(even+1)%3].Trxy
	// *wipe = make(map[int][3]int)
	// (*st).Trace[1].Unlock() ; (*st).Trace[2].Unlock() ; (*st).Trace[0].Unlock()
	read, write, epoch := st.WhichTrace()
	now := (epoch/tAxisStep)%tRange
	read.Lock() ; later := (*read).Trxy ; read.Unlock()
	write.Lock() ; trace := (*write).Trxy ; write.Unlock()
	traceLen := len(later)+len(trace)
	if traceLen == 0 { 
		write.Lock()
		renew := (*write).Trxy
		add := [4]int{ // race-6
			now,
			functions.ChancedRound( 2000*functions.Rand()-1000 )/250*250, 
			functions.ChancedRound( 20000*functions.Rand()-10000 ), 
			functions.ChancedRound( 20000*functions.Rand()-10000 ),
		}
		renew = append(renew, add)
		(*write).Trxy = renew
		write.Unlock()
		return 
	}
	latest, buffer := -functions.TRange, [][4]int{}
	for _, each := range later { 
		each[0] += -tRange
		buffer = append(buffer, each) 
	}// ; fmt.Println("READing old traces:", even, ts-tRange, each) }
	for _, each := range trace { buffer = append(buffer, each) }//; fmt.Println("READing current traces:", even, ts, each)}
	for ts, _ := range buffer { if ts > latest { latest = ts } }
	latestStep := buffer[latest]
	// (*st).Current.Lock()
	(*st.Current).Base.Lock() ; (*st.Current).Atts.Lock()
	distance := 1/math.Sqrt2 //(*st.Current.Atts).Agility // static yet
	id := functions.GetID( (*st.Current.Base).ID, (*st.Current.Atts).ID )
	(*st.Current).Base.Unlock() ; (*st.Current).Atts.Unlock()
	// (*st).Current.Unlock()
	angle := float64(latestStep[1])/1000 // * math.Pi / 180
	newAng := functions.Round((angle + rotate)*1000) // * math.Pi / 180
	for { if newAng > 1000 { newAng += -2000 } else if newAng < -1000 { newAng += 2000 } else { break }}
	newstep := [4]int{ now, newAng, latestStep[2], latestStep[3] }
	if step {
		turn := float64(newAng) / 1000 * math.Pi
		newstep[2] = functions.Round(float64(latestStep[2]) + 1000*distance*math.Sin(turn))
		newstep[3] = functions.Round(float64(latestStep[3]) + 1000*distance*math.Cos(turn))
	}
	toWrite := make(map[string][][3]int) // id: t, x, y
	for ts := latest ; ts < now ; ts++ { 
		toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[2], latestStep[3]})
	}
	write.Lock()
	rewrite := (*write).Trxy
	rewrite = append(rewrite, newstep) // race-2 race-6
	(*write).Trxy = rewrite
	write.Unlock()
	toWrite[id] = append(toWrite[id], [3]int{now, newstep[2], newstep[3]})
	writeToCache <- toWrite
	functions.Wait(float64(tAxisStep)*math.Pi)// / math.Log2(dTrace[1]ance+1)) // 1.536 - 0.256
}

func (st *State) Path() [5][2]int {
	tRange := functions.TRange
	// epoch := functions.Epoch()
	// even := (epoch/(tRange*tAxisStep))%3
	// (*st).Trace[1].Lock() ; (*st).Trace[2].Lock() ; (*st).Trace[0].Lock()
	// later, trace := (*st).Trace[(even+2)%3].Trxy, (*st).Trace[(even+3)%3].Trxy
	// (*st).Trace[1].Unlock() ; (*st).Trace[2].Unlock() ; (*st).Trace[0].Unlock()
	read, write, _ := st.WhichTrace()
	read.Lock() ; later := (*read).Trxy ; read.Unlock() // race-5 race-7
	write.Lock() ; trace := (*write).Trxy ; write.Unlock() // race-5 race-7
	if len(trace)+len(later) == 0 { return [5][2]int{} } // race-6
	buffer := make(map[int][4]int)
	for _, each := range later { buffer[each[0]-tRange] = each }
	for _, each := range trace { buffer[each[0]] = each }
	// for _, each := range trace { buffer = append(buffer, each) }
	tIndex, idx := [][4]int{}, 0
	for i:=tRange ; i>-tRange ; i-- {
		if _, ok := buffer[i] ; ok { tIndex = append(tIndex, buffer[i]) ; idx++ }
	// 	for j:=len(buffer)-1; j<0; j-- {
	// 		if buffer[j][0] == i { tIndex = append(tIndex, buffer[j]) ; idx++ ; break }
	// 	}
		if idx >= 4 { break }
	}
	// for i:=0 ; i<4 ; i++ {
	// 	idx := len(buffer)-1-i
	// 	if idx >= 0 { tIndex = append(tIndex, buffer[idx]) }
	// }
	if len(tIndex) == 0 { return [5][2]int{} }
	if len(tIndex) == 1 { 
		return [5][2]int{
			[2]int{ tIndex[0][1], 0 },
			[2]int{ tIndex[0][2], tIndex[0][3] },
			[2]int{ tIndex[0][2], tIndex[0][3] },
			[2]int{ tIndex[0][2], tIndex[0][3] },
			[2]int{ tIndex[0][2], tIndex[0][3] },
		}
	}
	if len(tIndex) == 2 {
		return [5][2]int{
			[2]int{ tIndex[0][1], -tIndex[0][1]+tIndex[1][1] },
			[2]int{ tIndex[0][2], tIndex[0][3] },
			[2]int{ tIndex[1][2], tIndex[1][3] },
			[2]int{ tIndex[1][2], tIndex[1][3] },
			[2]int{ tIndex[1][2], tIndex[1][3] },
		}	
	}
	if len(tIndex) == 3 {
		return [5][2]int{
			[2]int{ tIndex[0][1], -tIndex[0][1]+tIndex[1][1] },
			[2]int{ tIndex[0][2], tIndex[0][3] },
			[2]int{ tIndex[1][2], tIndex[1][3] },
			[2]int{ tIndex[2][2], tIndex[2][3] },
			[2]int{ tIndex[2][2], tIndex[2][3] },
		}	
	}
	return [5][2]int{
		[2]int{ tIndex[0][1], -tIndex[0][1]+tIndex[1][1] },
		[2]int{ tIndex[0][2], tIndex[0][3] },
		[2]int{ tIndex[1][2], tIndex[1][3] },
		[2]int{ tIndex[2][2], tIndex[2][3] },
		[2]int{ tIndex[3][2], tIndex[3][3] },
	}
}
