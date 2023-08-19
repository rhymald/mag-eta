package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play"
)

func around(c *gin.Context) { 
	base := play.Init_BasicStats()
	char := base.Init_Character()
	char.Init_Attributes()
	randomPath := [5][2]int{}
	randomPath[0][0] = 360
	randomPath[0][1] = -100
	c.IndentedJSON(200, struct{ Char *play.Character ; Simp play.Simplified }{Char: char, Simp: char.Simplify(randomPath)})
}
