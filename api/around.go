package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play"
)

func around(c *gin.Context) { 
	var buffer []play.Simplified
	takenID, myPlayer := "", &play.State{}
	if _, ok := c.Request.Header["myplayerid"] ; ok { takenID = c.GetHeader("myplayerid") } else { takenID = c.Param("myplayerid") }
	theWorld.ByID.Lock()
	if _, ok := (*theWorld.ByID).List[takenID] ; ok { myPlayer = (*theWorld.ByID).List[takenID] } else { myPlayer = nil }
	if myPlayer != nil { 
		path := myPlayer.Path()
		buffer = append(buffer, myPlayer.Current.Simplify( path )) 
	}
	for id, each := range theWorld.ByID.List {
		if id == takenID { continue }
		path := each.Path()
		buffer = append(buffer, each.Current.Simplify( path ))
	} 
	theWorld.ByID.Unlock()
	c.JSON(200, buffer)
}
