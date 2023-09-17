package api 

import (
	"github.com/gin-gonic/gin"
	// "rhymald/mag-eta/play"
	"strconv"
	"fmt"
)

func move(c *gin.Context) {
	takenID, direction := "", ""
	if _, ok := c.Request.Header["myplayerid"] ; ok { takenID = c.GetHeader("myplayerid") } else { takenID = c.Param("myplayerid") }
	if _, ok := c.Request.Header["direction"] ; ok { direction = c.GetHeader("direction") } else { direction = c.Param("direction") }
	if direction != "" && takenID != "" {
		where, err := strconv.Atoi(direction)
		theWorld.ByID.Lock()
		read := (*theWorld).ByID
		theWorld.ByID.Unlock()
		myPlayer, _ := read.Read(takenID) //; ok { myPlayer, _ = read.Read(takenID) } else { myPlayer = nil }
		if myPlayer != nil && err == nil {
			myPlayer.Current.Base.Lock()
			thread := theWorld.GimmeThread((*myPlayer.Current.Base).ID)
			myPlayer.Current.Base.Unlock()
			myPlayer.Move( float64(where)/1000, true, (*thread).Chan )
			c.IndentedJSON(200, "Moved")
		} else {
			c.IndentedJSON(400, []string{ "Bad request headers:", "- myplayerid parsed:", fmt.Sprint(myPlayer), "- direction parsed:", fmt.Sprint(err) })
		}
	} else {
		c.IndentedJSON(400, []string{ "Bad request headers:", "- myplayerid read:", takenID, "- direction read:", direction })
	}
}