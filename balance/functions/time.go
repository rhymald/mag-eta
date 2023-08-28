package functions

import (
	"time"
)

const (
	TAxisStep = 400 //ms for grid
	TRange = 1250 //steps per bucket
	LifeCycleQueuePause = 1618
)

var StartEpoch = int(time.Now().UnixNano())

func TAxis() int { return (Epoch()/TAxisStep)%TRange }
func Epoch() int { return EpochNS()/1000000 }
func EpochNS() int { return int(time.Now().UnixNano())-StartEpoch }
func Wait(ms float64) { time.Sleep( time.Millisecond * time.Duration( ms )) }
