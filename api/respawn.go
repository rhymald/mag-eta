package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play"
	"rhymald/mag-eta/balance/functions"
)

var theWorld = play.Init_World()

func login(c *gin.Context) { 
	base := play.Init_BasicStats()
	char := base.Init_Character()
	char.Init_Attributes()
	state := char.Init_State()
	id := theWorld.Login(state)
	_, writer := theWorld.WhichGrid()
	writer.Reg.Register(id)
	go func(){ state.Lifecycle_EffectConsumer() }()
	go func(){ state.Lifecycle_Regenerate() }()
	go func(){ for { state.Move( (1+functions.Rand())/8, true, theWorld.Queue.Chan ) }}()
	c.IndentedJSON(200, struct{ ID string ; Result string }{ ID: id, Result: "Successfully Logged In" })
}
