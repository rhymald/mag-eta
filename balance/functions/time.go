package functions

import (
	"time"
)

var StartEpoch = int(time.Now().UnixNano())

func Epoch() int { return EpochNS()/1000000 }
func EpochNS() int { return int(time.Now().UnixNano())-StartEpoch }
func Wait(ms float64) { time.Sleep( time.Millisecond * time.Duration( ms )) }
