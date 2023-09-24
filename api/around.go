package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play/character"
	"rhymald/mag-eta/balance/functions"
	"math"
	"log"
	"fmt"
)

func around(c *gin.Context) { 
	distanceLimit := 2000000000.0
	cutOff := functions.CeilRound(math.Log2(1+distanceLimit))
	objectLimit := 5 + 1
	var buffer []character.Simplified
	takenID, myPlayer := "", &character.State{}
	if _, ok := c.Request.Header["myplayerid"] ; ok { takenID = c.GetHeader("myplayerid") } else { takenID = c.Param("myplayerid") }
	theWorld.ByID.Lock()
	read, filter := (*theWorld).ByID, (*theWorld).Grid
	theWorld.ByID.Unlock()
	myPlayer, _ = read.Read(takenID) //; ok { myPlayer, _ = read.Read(takenID) } else { myPlayer = nil }
	first := [2]int{}
	objectLimit += functions.CeilRound( math.Sqrt( float64(read.Len())+1 ))
	if myPlayer != nil { path := myPlayer.Path() ; first = path[1] }
	fmt.Println(">>--> [0] Center:", first, "- has ID:", takenID)
	// buffer = append(buffer, myPlayer.Current.Simplify( path )) 
	counter := 0
	okies, passen, notokies := 0.0, 0.0, 0.0
	allstates := read.GetAll()
	filtered := filter.GetAgainst(0, first)//[:cutOff]
	fmt.Println(">>--> [0] filtered len:", len(filtered), "- and distance limit:", cutOff)
	if len(filtered) > cutOff {  filtered = filtered[:cutOff] }
	for _, listed := range filtered { for _, id := range listed {
		each, ok := allstates[id]
		if ok {
			fmt.Println(">>--> [1] ID found:", id)
			path := each.Path()
			// if path == [5][2]int{} { log.Fatalln(functions.FatalErrors["NotEnoughCPU"]) }
			if path == [5][2]int{} { notokies += 1 ; continue } else { okies += 1 }
			// beyond := functions.Vector( float64(path[1][0]-first[0]), float64(path[1][1]-first[1]) ) > distanceLimit
			if counter >= objectLimit { continue } // replace with break
			buffer = append(buffer, each.Current.Simplify( path ))
			fmt.Println(">>----> [2] ID add:", id, path)
			counter++	
		} else { fmt.Println(">>--> [1] No such ID:", id)	}
	}}
	// for id, each := range allstates { // race-3 race-8 }
	if notokies != 0 { log.Fatalf("%v - Okies: %0.0f, Stuck: %0.0f, Passen: %0.0f", functions.PanicErrors["NotEnoughCPU"], okies, notokies, passen) }
	c.JSON(200, buffer)
}
