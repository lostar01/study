package router

import (
	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/apis"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", apis.Ping)
	r.GET("/namespace", apis.GetNameSpace)
	r.GET("/deployment/list/:namespace", apis.ListDeployment)
}
