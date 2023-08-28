package api

import (
	"github.com/gin-gonic/gin"
	"rhymald/mag-eta/balance/functions"
	"rhymald/mag-eta/play"
	"math"
	// "fmt"
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
	list := theWorld.Seek_Square(0, 0, 1414)
	// reader, writer := theWorld.WhichGrid()
	// buffer := writer.Get_Square(0, 0, 1414)
	// old := reader.Get_Square(0, 0, 1414)
	// for id, row := range old { if _, ok := buffer[id] ; !ok { 
	// 	row[2] += -functions.TRange
	// 	buffer[id] = row
	// }}
	// list := []string{}
	// states := (*theWorld).ByID.GetAll()
	// for t:=functions.TRange-1 ; t>=-functions.TRange ; t-- { for id, row := range buffer {
	// 	player, ok := states[id]
	// 	path := player.Path()
	// 	actual := row[0] == path[1][0] && row[1] == path[1][1] && ok
	// 	if row[2] == t && actual {
	// 		list = append(list, fmt.Sprintf("id = %s, x = %6d, y = %6d, t = %3d, old = %6dms", id, row[0], row[1], row[2], row[3]))
	// 	}
	// }}
	// buffer.PosR[0], buffer.PosR[1] = reader.Get_CentralPos()
	// buffer.PosW[0], buffer.PosW[1] = writer.Get_CentralPos()
	// reader.Lock()
	// buffer.Xr = reader.X.Get(targetX, targetAOE, targetT)
	// buffer.Yr = reader.Y.Get(targetY, targetAOE, targetT)
	// reader.Unlock()
	// writer.Lock()
	// buffer.Xw = writer.X.Get(targetX, targetAOE, targetT)
	// buffer.Yw = writer.Y.Get(targetY, targetAOE, targetT)
	// writer.Unlock()
	c.IndentedJSON(200, list)
}

func spawn(c *gin.Context) { 
	base := play.Init_BasicStats()
	char := base.Init_Character()
	char.Init_Attributes()
	state := char.Init_State()
	id := theWorld.Login(state)
	go func(){ functions.Wait(100) ; state.Lifecycle_EffectConsumer() }()
	go func(){ functions.Wait(600) ; state.Lifecycle_Regenerate() }()
	direction := functions.Rand() - functions.Rand()
	direction = direction / (math.Abs(direction))
	go func(){ functions.Wait(4000) ; for { state.Move( direction*(1+2*functions.Rand())/24, true, theWorld.Queue.Chan ) }}()
	c.IndentedJSON(200, struct{ ID string ; Result string }{ ID: id, Result: "Successfully spawned" })
}
