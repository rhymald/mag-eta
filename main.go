package main

import (
	"rhymald/mag-eta/api"
	// "rhymald/mag-eta/balance/functions"
	"fmt"
)

var (
	ipAddress = ":4917"
)

func init() {
	fmt.Println(" >>>>> Application is getting started... <<<<<")
}

func main() {
	api.Init_API(ipAddress)
}
