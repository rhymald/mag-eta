package main

import (
	"rhymald/mag-eta/api"
	"rhymald/mag-eta/balance/functions"
	"flag"
	"fmt"
)

var (
	ipAddress = ":4917"
)

func init() {
	flag.IntVar(&functions.Threads, "t", 4, "Threads to execute global queues")
	fmt.Println(" >>>>> Application is getting started... <<<<<")
}

func main() {
	api.Init_API(ipAddress)
}