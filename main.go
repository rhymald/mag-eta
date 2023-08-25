package main

import (
	"rhymald/mag-eta/api"
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