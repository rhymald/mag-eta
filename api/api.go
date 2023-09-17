package api

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"fmt"
)

type LogEvent struct {
	Source  string `json:"Source"`
	Method  string `json:"Method"`
	Path    string `json:"Path"`
	Status  int    `json:"Status"`
	Latency string `json:"Latency"`
	Time    string `json:"Time"`
}

func Init_API(ipAddr string) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(jsonLoggerMiddleware())
	router.GET("/", around)
	router.GET("/:myplayerid", around)
	router.GET("/login", login)
	router.GET("/test/world", testWorld)
	router.GET("/test/spawn", spawn)
	router.GET("/move/:myplayerid/:direction", move)
	fmt.Println(" >>>>>      Hello, artifical world!      <<<<<")
	router.Run(ipAddr)
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			if params.StatusCode == 200 && params.Path == "/" { return "" }
			log := LogEvent{
				Status:  params.StatusCode,
				Latency: fmt.Sprintf("%0.3fms", float64(params.Latency.Microseconds())/1000),
				Method:  params.Method,
				Path:    params.Path,
				Time:    params.TimeStamp.Format("2006/01/02 15:04:05.999"),
				Source:  params.ClientIP,
			}
 			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}