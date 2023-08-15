package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play"
)

func around(c *gin.Context) { 
	base := play.Init_BasicStats()
	char := base.Init_Character()
	char.Init_Attributes()
	c.IndentedJSON(200, char)
}
