package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/balance/functions"
)

type Testing_Response struct {
	PosW [2]int
	Xw map[string][3]int
	Yw map[string][3]int
	PosR [2]int
	Xr map[string][3]int
	Yr map[string][3]int
}

func testWorld(c *gin.Context) {
	targetX, targetY, targetAOE, targetT := 0, 0, 2000, functions.TAxis()
	var buffer Testing_Response
	reader, writer := theWorld.WhichGrid()
	buffer.PosR[0], buffer.PosR[1] = reader.Get_CentralPos()
	buffer.PosW[0], buffer.PosW[1] = writer.Get_CentralPos()
	reader.Lock()
	buffer.Xr = reader.X.Get(targetX, targetAOE, targetT)
	buffer.Yr = reader.X.Get(targetY, targetAOE, targetT)
	reader.Unlock()
	writer.Lock()
	buffer.Xw = writer.X.Get(targetX, targetAOE, targetT)
	buffer.Yw = writer.X.Get(targetY, targetAOE, targetT)
	writer.Unlock()
	c.IndentedJSON(200, buffer) // world too big to output
}
