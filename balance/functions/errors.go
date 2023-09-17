package functions

import (
	"errors"
	"runtime"
)

var Threads = runtime.GOMAXPROCS(runtime.NumCPU()-1)

type ErrorList map[string]error 
var (
	FatalErrors = ErrorList{
		"NotEnoughCPU": errors.New("Not enough resources to handle characters and objects."),
		"NoThreadAssigned": errors.New("Location has no queue for this thread."),
	}
)