package functions

import (
	"errors"
)

type ErrorList map[string]error 

var (
	FatalErrors = ErrorList{
		"NotEnoughCPU": errors.New("Not enough resources to handle characters and objects."),
	}
)