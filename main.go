package main

import (
	"rhymald/mag-eta/api"
	// "rhymald/mag-eta/balance/functions"
	"runtime"
	"fmt"
)

var (
	ipAddress = ":4917"
)

func init() {
	threads := runtime.NumCPU()
	used := runtime.GOMAXPROCS( threads )
	fmt.Println(" >>---> Application is getting started...")
	fmt.Printf(" >>-----> Running on %d threads, total %d threads.\n", used, threads)
}
func main() {
	// go func(){ for {
	// 	functions.Wait(10000) 
	// 	fmt.Println(" >>---> Goroutines are running:", runtime.NumGoroutine())
	// }}()
	fmt.Println(" >>-------> Hello, artifical world! API is active.")
	api.Init_API(ipAddress)
}