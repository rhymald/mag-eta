package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play"
)

var theWorld = play.Init_World()

func login(c *gin.Context) { 
	base := play.Init_BasicStats()
	char := base.Init_Character()
	char.Init_Attributes()
	state := char.Init_State()
	id := theWorld.Login(state)
	go func(){ state.Lifecycle_EffectConsumer() }()
	go func(){ state.Lifecycle_Regenerate() }()
	c.IndentedJSON(200, struct{ ID string ; Result string }{ ID: id, Result: "Successfully Logged In" })
}
