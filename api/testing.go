package api

import (
	"github.com/gin-gonic/gin"
	// "rhymald/mag-eta/balance/functions"
	"rhymald/mag-eta/play/character"
	"rhymald/mag-eta/balance/primitives"
	// "math"
	// "fmt"
)

// type Testing_Response struct {
// 	PosW [2]int
// 	Xw map[string][3]int
// 	Yw map[string][3]int
// 	PosR [2]int
// 	Xr map[string][3]int
// 	Yr map[string][3]int
// }

func testWorld(c *gin.Context) {
	// list := theWorld.Seek_Square(0, 0, 1414)
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
	output := theWorld.Grid.GetAgainst(3)
	c.IndentedJSON(200, output)
}

func spawn(c *gin.Context) { 
	base := character.Init_BasicStats(primitives.PhysList[2])
	char := base.Init_Character()
	char.Init_Attributes()
	state := char.Init_State()
	state.Current.Base.Lock()
	thread := theWorld.GimmeThread((*state.Current.Base).ID)
	state.Current.Base.Unlock()
	// myPlayer.Move( float64(where)/1000, true, (*thread).Chan )
	state.Move( 0, false, (*thread).Chan )
	id := theWorld.Login(state)
	go func(){ state.Lifecycle_EffectConsumer() }()
	go func(){ state.Lifecycle_Regenerate() }()
	// direction := functions.Rand() - functions.Rand()
	// direction = direction / (math.Abs(direction))
	go func(){ for { state.Move( 1/24, true, (*thread).Chan ) }}()
	c.IndentedJSON(201, struct{ ID string ; Result string }{ ID: id, Result: "Successfully spawned" })
}
