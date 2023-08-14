package api

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"fmt"
)

type LogEvent struct {
	Status  int    `json:"Status"`
	Time    string `json:"Time"`
	Latency string `json:"Latency"`
	Source  string `json:"Source"`
	Method  string `json:"Method"`
	Path    string `json:"Path"`
}

func RunAPI(ipAddr string) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(jsonLoggerMiddleware())
	router.GET("/", hiThere)
	router.Run(ipAddr)
}

func hiThere(c *gin.Context) { 
	c.IndentedJSON(200, "Hello world!")
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			if len(params.Path) >= 7 {
				if params.StatusCode == 200 && params.Path[:7] == "/around" { return "" }
			}
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