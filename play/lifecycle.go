package play

import (
	"rhymald/mag-eta/balance/functions"
	"rhymald/mag-eta/balance/primitives"
	"math"
	"fmt"
)

func (st *State) Lifecycle_Regenerate() {
	isNPC := false //(*st).Current.IsNPC()
	if !isNPC { for {
		(*st).Current.Lock()
		energyFull := len((*st.Current.Pool).Dots) >= functions.FloorRound((*st.Current.Atts).Capacity)
		count, mean := len((*st.Current.Base).Energy), 0.0
		if count == 0 || energyFull { (*st).Current.Unlock() ; functions.Wait(4236) ; continue }
		(*st).Current.Unlock()
		// _, span := tracer.Start(ctx, fmt.Sprintf("regen-player-%d", (*st).Current.GetID()))
		// defer span.End()
		effect := primitives.Init_Effect()
		(*st).Current.Lock()
		if (*st).Current.HP.Dead() { (*st).Current.Unlock() ; return } // span.AddEvent("Player died") ; return }
		for picker:=0 ; picker<count ; picker++ {
			stream := (*st.Current.Base).Energy[picker]
			dot := stream.Init_Dot()
			pause := effect.Add_Self_MakeDot(dot)
			mean += 1/pause
		}
		wait := float64(count) / mean
		(*st).Current.Unlock()
		effect.Add_Self_HPRegen(16)
		st.Lock()
		renew := (*st).Effects
		renew[(*effect).Time] = effect
		(*st).Effects = renew
		// span.AddEvent(fmt.Sprintf("Emmiting regeneration effect: { %+v }", *effect))
		st.Unlock()
		functions.Wait(wait/4)
	}} else { for {
		// _, span := tracer.Start(ctx, fmt.Sprintf("regen-npc-%d", (*st).Current.GetID()))
		// defer span.End()
		if (*st.Current).HP.Wounded() { return } //span.AddEvent("NPC died") ; return }
		effect := primitives.Init_Effect()
		effect.Add_Self_HPRegen(64)
		st.Lock()
		renew := (*st).Effects
		renew[(*effect).Time] = effect
		(*st).Effects = renew
		st.Unlock()
		// span.AddEvent(fmt.Sprintf("Emmiting regeneration effect: { %+v }", *effect))
		functions.Wait(4236)
	}}
}

func (st *State) Lifecycle_EffectConsumer() {
	pause := functions.LifeCycleQueuePause
	// prefix := "player" ; if (*st).Current.IsNPC() { prefix = "npc" }
	for {
		first, sum, counter := 0, 0, 0
		now := functions.Epoch()
		// _, span := tracer.Start(ctx, fmt.Sprintf("lifecycle-%s-%d", prefix, (*st).Current.GetID()))
		// defer span.End()
		st.Lock() 
		startLen := len((*st).Effects) 
		sleep := float64(pause) / math.Pow( math.Phi, math.Log2(2+float64(startLen))/math.Log2(3)-1 )
		st.Unlock()
		if startLen == 0 { functions.Wait(float64(pause+1000)) ; continue }
		st.Lock() ; fmt.Println("BROKER[0] Before:", len((*st).Effects), sum) ; st.Unlock()

		// step 1 read to limit
		buffer := make(map[int]*primitives.Effect)
		st.Lock()
		// startLen := len((*st).Effects)
		for ts, effect := range (*st).Effects {
			if len(buffer) == 0 { first = ts }
			if ts-first > pause { continue }
			buffer[ts] = effect
			counter++ ; sum += ts - first
			if sum > counter * pause { break }
		}
		// span.AddEvent(fmt.Sprintf("Effects: { read: %d, total: %d }", counter, startLen))
		st.Unlock()

		// TBD conditions
		instant, _, delayed := make(map[int]primitives.Different_Effects), make(map[int]primitives.Different_Effects), make(map[int]primitives.Different_Effects)
		counterInst, _, counterDel := 0, 0, 0
		// step 2 sum and distribute
		for ts, each_effect := range buffer { for idx, each := range (*each_effect).Effects {
			switch kind := each.(type) {
			case primitives.Effect_HPRegen:
				tsNew := ts
				for { if _, ok := instant[tsNew]; ok { tsNew = tsNew+1 } else {break} }
				instant[tsNew] = each
				counterInst++
			case primitives.Effect_MakeDot:
				tsNew := (ts + each.Delayed()) - now
				for { if _, ok := delayed[tsNew]; ok { tsNew = tsNew+1 } else {break} }
				delayed[tsNew] = each
				counterDel++
			default:
				fmt.Printf("Unknown sub-effect[%d][%d] type[%v]: %+v\n", ts, idx, kind, each)
				// span.RecordError(errors.New(fmt.Sprintf("Unknown sub-effect[%d][%d] type[%v]: %+v", ts, idx, kind, each)))
			}
		}}
		// span.AddEvent(fmt.Sprintf("Sorted: { instant: %d, delayed: %d, conditions: 0, total: %d }", counterInst, counterDel, counterDel+counterInst))
		
		// step 3 cut condies

		// step 4 redirect back leftovers
		threshold := functions.CeilRound(sleep / math.Phi)
		accumulator := -threshold
		counterDelayed, counterTransformed := []string{}, []string{}
		for diff, each := range delayed {
			if accumulator + diff < pause - threshold {
				tsNew := now + diff
				for { if _, ok := instant[tsNew]; ok { tsNew = tsNew+1 } else {break} }
				instant[tsNew] = each
				delete(delayed, diff)
				counterTransformed = append(counterTransformed, fmt.Sprintf("%+d", diff))
				// span.AddEvent(fmt.Sprintf("Saved for consume: { %+v }", each))
			} else {
				tsNew := now - diff - diff / 7
				st.Lock()
				renew := (*st).Effects
				for { if _, ok := renew[tsNew]; ok { tsNew = tsNew+1 } else {break} }
				sentBack := primitives.Init_Effect()
				(*sentBack).Time = tsNew
				(*sentBack).Effects = append((*sentBack).Effects, each)
				renew[tsNew] = sentBack
				(*st).Effects = renew
				st.Unlock()
				// span.AddEvent(fmt.Sprintf("Redirected back to queue: { %+v }", *sentBack))
				counterDelayed = append(counterDelayed, fmt.Sprintf("%+d", diff))
			}
			accumulator += diff
		}
		fmt.Println("BROKER[4] Filtered:")
		fmt.Println("  Back to queue:", counterDelayed)
		fmt.Println("  Consumed:     ", counterTransformed)
		
		// step 5 consume instants
		hpregens := 0
		makedots := make(map[int]primitives.Dot)
		for time, each := range instant {
			switch kind := each.(type) {
			case primitives.Effect_HPRegen:
				hpregens += each.HP()
			case primitives.Effect_MakeDot:
				ts, dots := each.Dots()
				for _, dot := range dots { makedots[ts+time] = dot }
			default:
				fmt.Printf("Unknown sub-effect type[%v]: %+v\n", kind, each)
				// span.RecordError(errors.New(fmt.Sprintf("Unknown sub-effect type[%v] to apply: %+v", kind, each)))
			}
		}
		fmt.Println("BROKER[5] To apply:")
		fmt.Println("  HP:", hpregens)
		fmt.Println("  Dots:", makedots)
		(*st).Current.Lock()
		(*st.Current).HP.Upd(hpregens)
		// (*(*st).Current).ID["Life"] = functions.Epoch()
		// span.AddEvent(fmt.Sprintf("HP modificated: %+d", hpregens))
		for ts, dot := range makedots {
			tsNew := ts
			(*st.Current).Pool.Lock()
			renew := (*st.Current.Pool).Dots
			for { if _, ok := renew[tsNew]; ok { tsNew = tsNew+1 } else { break } }
			renew[tsNew] = &dot
			(*st.Current.Pool).Dots = renew
			(*st.Current).Pool.Unlock()
		}
		// span.AddEvent(fmt.Sprintf("Dots to append: { %v }", makedots))
		(*st).Current.Unlock()
		fmt.Println("BROKER[5] Applied")

		// step F clean read from queue
		st.Lock() 
		for ts, _ := range buffer { delete((*st).Effects, ts) }
		// span.AddEvent(fmt.Sprintf("Effects: { read: %d, total before: %d, total after: %d }", counter, startLen, len((*st).Effects)))
		fmt.Println("BROKER[F] After:", len((*st).Effects), sum)
		st.Unlock()
		// span.AddEvent(fmt.Sprintf("Wait for: %0.3fms", sleep ))
		// span.End()
		fmt.Println("BROKER[F] Sleep:", sleep, "ms")
		fmt.Println("--------------------------------------------")
		functions.Wait( sleep )
	}
}