package functions

import (
	"errors"
)

type ErrorList map[string]error 

var (
	PanicErrors = ErrorList{
		"NotEnoughCPU":         errors.New("Not enough resources to handle characters and objects."),
		"UnexpectedTargetType": errors.New("Target position is neither coordinates or ID."),
	}
	FatalErrors = ErrorList{
		"NoThreadAssigned": errors.New("Location has no queue for this thread."),
	}
)