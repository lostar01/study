package middleware

import "github.com/gin-gonic/gin"

func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-COntrol-Allow-Headers", "authorizatioin,origin,content-type,accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIOINS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}
