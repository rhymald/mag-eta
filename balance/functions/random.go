package functions

import (
	"math"
	"time"
	"math/rand"
	"crypto/sha512"
	"encoding/binary"	
)

const MinEntropy = 0.0132437

func Rand() float64 {
  x := (time.Now().UnixNano())
  in_bytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(in_bytes, uint64(x))
  hsum := sha512.Sum512(in_bytes)
  sum  := binary.BigEndian.Uint64(hsum[:])
  return rand.New(rand.NewSource( int64(sum) )).Float64()
}

func Ntrp(a float64) float64 { 
  randy := (Epoch() % 1000) / 333
  entropy := math.Log10( math.Abs(a)+1 )/25 
  if randy == 2 { a = a*(1+MinEntropy+entropy) }
  if randy == 0 { a = a/(1+MinEntropy+entropy) }
  return math.Round( a*1000 ) / 1000
}

func Vector(args ...float64) float64 { 
	sum := 0.0
	for _,each := range args { sum += each*each }
	return math.Sqrt(sum)
}
