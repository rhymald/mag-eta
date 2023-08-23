package api

import (
	"github.com/gin-gonic/gin"
	// "rhymald/mag-eta/play"
)

func testWorld(c *gin.Context) {
	c.IndentedJSON(200, theWorld.Grid) // world too big to output
}
