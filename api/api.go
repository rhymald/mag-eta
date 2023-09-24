package api

import (
	"github.com/gin-gonic/gin"
	"runtime"
	"encoding/json"
	"fmt"
)

type LogEvent struct {
	// Method  string `json:"Method"`
	Time    string `json:"Time"`
	Path    string `json:"Path"`
	Status  int    `json:"Status"`
	Latency string `json:"Latency"`
	Source  string `json:"Source"`
	GRpO    string `json:"GRpO"`
	// Query   map[string]string `json:"Query"`
	// Metrics map[string]int    `json:"Metrics"`
}

func Init_API(ipAddr string) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(jsonLoggerMiddleware())
	router.GET("/", around)
	router.GET("/at/:myplayerid", around)
	router.GET("/login", login)
	router.GET("/test/world", testWorld)
	router.GET("/test/spawn", spawn)
	router.GET("/move/:myplayerid/:direction", move)
	router.Run(ipAddr)
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			if params.StatusCode == 200 && params.Path == "/" { return "" }
			routines := runtime.NumGoroutine()
			objects := theWorld.ByID.Len()
			log := LogEvent{
				Latency: fmt.Sprintf("%0.3fms", float64(params.Latency.Microseconds())/1000),
				Time:    params.TimeStamp.Format("2006/01/02 15:04:05"),
				Path:    fmt.Sprintf("%s[%s]", params.Method, params.Path),
				GRpO:    fmt.Sprintf("%d/%d", routines, objects),
				Status:  params.StatusCode,
				Source:  params.ClientIP,
				// Method:  params.Method,
				// Objects: theWorld.ByID.Len(),
				// Routines: runtime.NumGoroutine(),
				// Metrics: make(map[string]int),
				// Query:   make(map[string]int),
			}
			// log.Metrics["WorldObjects"] = theWorld.ByID.Len()
			// log.Metrics["GoRoutines"] = runtime.NumGoroutine()
 			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}