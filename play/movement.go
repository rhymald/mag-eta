package play 

import (
	"rhymald/mag-eta/balance/functions"
	"math"
)

func (st *State) Move(rotate float64, step bool, writeToCache chan map[string][][3]int) {
	tAxisStep, tRange := functions.TAxisStep, functions.TRange
	epoch := functions.Epoch()
	now, even := (epoch/tAxisStep)%tRange, (epoch/(tRange*tAxisStep))%3
	(*st).Trace[0].Lock() ; (*st).Trace[1].Lock() ; (*st).Trace[2].Lock()
	later, trace, wipe := (*st).Trace[(even+2)%3].Trxy, (*st).Trace[(even+3)%3].Trxy, &(*st).Trace[(even+1)%3].Trxy
	*wipe = make(map[int][3]int)
	traceLen := len(later)+len(trace)
	(*st).Trace[1].Unlock() ; (*st).Trace[2].Unlock() ; (*st).Trace[0].Unlock()
	if traceLen == 0 { 
		(*st).Trace[even].Lock()
		renew := (*st).Trace[even].Trxy
		renew[now] = [3]int{ 
			functions.ChancedRound( 2000*functions.Rand()-1000 )/250*250, 
			functions.ChancedRound( 2000*functions.Rand()-1000 ), 
			functions.ChancedRound( 2000*functions.Rand()-1000 ),
		}
		(*st).Trace[even].Trxy = renew
		(*st).Trace[even].Unlock()
		return 
	}
	latest, buffer := -tRange, make(map[int][3]int)
	for ts, each := range later { buffer[ts-tRange] = each }// ; fmt.Println("READing old traces:", even, ts-tRange, each) }
	for ts, each := range trace { buffer[ts] = each }//; fmt.Println("READing current traces:", even, ts, each)}
	for ts, _ := range buffer { if ts > latest { latest = ts } }
	latestStep := buffer[latest]
	(*st).Current.Lock()
	distance := 1/math.Sqrt2 //(*st.Current.Atts).Agility // static yet
	id := functions.GetID( (*st.Current.Base).ID, (*st.Current.Atts).ID )
	(*st).Current.Unlock()
	angle := float64(latestStep[0])/1000 // * math.Pi / 180
	newAng := functions.Round((angle + rotate)*1000) // * math.Pi / 180
	for { if newAng > 1000 { newAng += -2000 } else if newAng < -1000 { newAng += 2000 } else { break }}
	newstep := [3]int{ newAng, latestStep[1], latestStep[2] }
	if step {
		turn := float64(newAng) / 1000 * math.Pi
		newstep[1] = functions.Round(float64(latestStep[1]) + 1000*distance*math.Sin(turn))
		newstep[2] = functions.Round(float64(latestStep[2]) + 1000*distance*math.Cos(turn))
	}
	toWrite := make(map[string][][3]int) // id: t, x, y
	for ts := latest ; ts < now ; ts++ { 
		toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
	}
	(*st).Trace[even].Lock()
	renew := (*st).Trace[even].Trxy
	renew[now] = newstep // here-2
	(*st).Trace[even].Trxy = renew
	(*st).Trace[even].Unlock()
	toWrite[id] = append(toWrite[id], [3]int{now, newstep[1], newstep[2]})
	writeToCache <- toWrite
	functions.Wait(float64(tAxisStep)*math.Pi)// / math.Log2(dTrace[1]ance+1)) // 1.536 - 0.256
}

func (st *State) Path() [5][2]int {
	tAxisStep, tRange := functions.TAxisStep, functions.TRange
	epoch := functions.Epoch()
	even := (epoch/(tRange*tAxisStep))%3
	(*st).Trace[1].Lock() ; (*st).Trace[2].Lock() ; (*st).Trace[0].Lock()
	later, trace := (*st).Trace[(even+2)%3].Trxy, (*st).Trace[(even+3)%3].Trxy
	(*st).Trace[1].Unlock() ; (*st).Trace[2].Unlock() ; (*st).Trace[0].Unlock()
	if len(trace)+len(later) == 0 { return [5][2]int{} }
	buffer := make(map[int][3]int)
	for ts, each := range later { buffer[ts-tRange] = each }
	for ts, each := range trace { buffer[ts] = each } // here-2
	tIndex, idx := []int{}, 0
	for i:=tRange ; i>-tRange ; i-- {
		if _, ok := buffer[i] ; ok { tIndex = append(tIndex, i) ; idx++ }
		if idx >= 4 { break }
	}
	if idx <= 1 { 
		return [5][2]int{
			[2]int{ buffer[tIndex[0]][0], 0 },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
		}
	}
	if idx == 2 {
		return [5][2]int{
			[2]int{ buffer[tIndex[0]][0], -buffer[tIndex[0]][0]+buffer[tIndex[1]][0] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
		}	
	}
	if idx == 3 {
		return [5][2]int{
			[2]int{ buffer[tIndex[0]][0], -buffer[tIndex[0]][0]+buffer[tIndex[1]][0] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
			[2]int{ buffer[tIndex[2]][1], buffer[tIndex[2]][2] },
			[2]int{ buffer[tIndex[2]][1], buffer[tIndex[2]][2] },
		}	
	}
	return [5][2]int{
		[2]int{ buffer[tIndex[0]][0], -buffer[tIndex[0]][0]+buffer[tIndex[1]][0] },
		[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
		[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
		[2]int{ buffer[tIndex[2]][1], buffer[tIndex[2]][2] },
		[2]int{ buffer[tIndex[3]][1], buffer[tIndex[3]][2] },
	}
}
