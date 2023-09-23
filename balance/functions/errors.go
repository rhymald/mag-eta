package functions

import (
	"errors"
)

type ErrorList map[string]error 

var (
	Warnings = ErrorList{
		"NotEnoughCPU": errors.New("Not enough resources to handle characters and objects."),
	}
	FatalErrors = ErrorList{
		"NoThreadAssigned": errors.New("Location has no queue for this thread."),
	}
)
