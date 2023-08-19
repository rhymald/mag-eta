package functions

import (
	"math"
)

func Round(a float64) int { return int(math.Round(a)) }
func CeilRound(a float64) int { return int(math.Ceil(a)) } 
// func FloorRound(a float64) int { return int(math.Floor(a)) }
func ChancedRound(a float64) int {
	b, l := math.Ceil(a), math.Floor(a)
	c := math.Abs(math.Abs(a)-math.Abs(math.Min(b, l)))
	if a<0 {c = 1-c}
	if Rand() < c {return int(b)} else {return int(l)}
	return 0
}