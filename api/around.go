package api

import (
	"github.com/gin-gonic/gin"
)

func around(c *gin.Context) { 
	c.IndentedJSON(200, "Hello world!")
}
