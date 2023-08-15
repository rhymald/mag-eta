package primitives 

import (
	"rhymald/mag-eta/balance/functions"
)

type Stream struct {
	E string
	C int 
	A int
	D int 
}

func Init_Stream(elem string) *Stream {
	maj, min, dev := 0.0, 0.0, 0.0
	for x:=0; x<4; x++ { maj += functions.Rand() ; min += functions.Rand() ; dev += functions.Rand() }
	leng := functions.Vector(maj, min, dev) 
	return &Stream{ 
		E: elem, 
		C: functions.CeilRound( maj/leng * 1000 ), 
		A: functions.CeilRound( min/leng * 1000 ), 
		D: functions.CeilRound( dev/leng * 1000 ), 
	}
}

func (str *Stream) Mean() float64 {
	return 3 / (1000/float64((*str).C) + 1000/float64((*str).A) + 1000/float64((*str).D))
}