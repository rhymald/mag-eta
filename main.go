package main

import (
	"rhymald/mag-eta/api"
)

var (
	ipAddress = ":4917"
)

func init() {}
func main() {
	api.Init_API(ipAddress)
}