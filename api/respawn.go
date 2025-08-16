package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/play/character"
	"rhymald/mag-eta/play/world"
	// "rhymald/mag-eta/balance/functions"
	"rhymald/mag-eta/balance/primitives"
)

var theWorld = world.Init_World()
var AllowLogin = true

func login(c *gin.Context) { 
	base := character.Init_BasicStats(primitives.PhysList[1])
	char := base.Init_Character()
	char.Init_Attributes()
	state := char.Init_State() ; state.Move( 0, false, theWorld.Queue.Chan ) ; state.Move( 0, false, theWorld.Queue.Chan )
	id := theWorld.Login(state) // does not work with empty trace
	// _, writer := theWorld.WhichGrid()
	// writer.Reg.Register(id)
	go func(){ state.Lifecycle_EffectConsumer() }()
	go func(){ state.Lifecycle_Regenerate() }()
	// direction := functions.Rand() - functions.Rand()
	// direction = direction / (math.Abs(direction))
	// go func(){ functions.Wait(4000) ; for { state.Move( direction*(1+2*functions.Rand())/24, true, theWorld.Queue.Chan ) }}()
	c.IndentedJSON(201, struct{ ID string ; Result string }{ ID: id, Result: "Successfully Logged In" })
}
