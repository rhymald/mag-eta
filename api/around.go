package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play"
	"rhymald/mag-eta/balance/functions"
	"math"
)

func around(c *gin.Context) { 
	distanceLimit := 10000.0
	objectLimit := 5 + 1
	var buffer []play.Simplified
	takenID, myPlayer := "", &play.State{}
	if _, ok := c.Request.Header["myplayerid"] ; ok { takenID = c.GetHeader("myplayerid") } else { takenID = c.Param("myplayerid") }
	theWorld.ByID.Lock()
	read := (*theWorld.ByID).List
	theWorld.ByID.Unlock()
	if _, ok := read[takenID] ; ok { myPlayer = read[takenID] } else { myPlayer = nil }
	first := [2]int{}
	objectLimit += functions.CeilRound( math.Sqrt( float64(len(read))+1 ))
	if myPlayer != nil { 
		path := myPlayer.Path() ; first = path[1]
		buffer = append(buffer, myPlayer.Current.Simplify( path )) 
	}
	counter := 0
	for id, each := range read {
		path := each.Path()
		beyond := functions.Vector( float64(path[1][0]-first[0]), float64(path[1][1]-first[1]) ) > distanceLimit
		if id == takenID || counter >= objectLimit || beyond { continue }
		buffer = append(buffer, each.Current.Simplify( path ))
		counter++
	} 
	c.JSON(200, buffer)
}
