package primitives 

import (
	"sync"
)

type Health struct {
	In  int
	Out int
	// BIn
	// BOut
	// BOff
	sync.Mutex
}

func Init_Health() *Health { return &Health{ In: 618, Out: 382 }}

func (hp *Health) Is() int { 
	hp.Lock()
	current := (*hp).In - (*hp).Out
	hp.Unlock()
	return current
}

func (hp *Health) Wounded() bool { if hp.Is() < 0 { return true } ; return false }
func (hp *Health) Full() bool { if hp.Is() >= 1000 { return true } ; return false }
func (hp *Health) Dead() bool { if hp.Is() <= -1000 { return true } ; return false }

func (hp *Health) Upd(amount int) { 
	if hp.Dead() || hp.Full() { amount = 0 }
	if hp.Wounded() { if amount > 0 { amount = 1 } else { amount = -1 }}
	if amount >= 1000 { amount = 999 } else if amount <= -1000 { amount = -999 } 
	if amount + hp.Is() > 1000 { amount = 1000 - hp.Is() }
	if amount > 0 { hp.Lock() ; (*hp).In += amount ; hp.Unlock() }
	if amount < 0 { hp.Lock() ; (*hp).Out += amount ; hp.Unlock() }
}
