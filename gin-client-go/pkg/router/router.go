package router

import (
	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/apis"
	"lostar.com/m/pkg/middleware"
)

func InitRouter(r *gin.Engine) {
	middleware.InitMiddleware(r)
	r.GET("/ping", apis.Ping)
	r.GET("/namespace/list", apis.GetNameSpace)
	r.GET("/deployment/list/:namespace", apis.ListDeployment)
	r.GET("/deployment/get/:namespace/:name", apis.GetDeployment)
	r.GET("/node/list", apis.GetNodeList)
	r.GET("/pod/list/:namespace", apis.GetPodList)
	r.GET("/pod/get/:namespace/:name", apis.GetPod)
	r.GET("/pod/exec/:namespaceName/:podName/:containerName", apis.ExecContainer)
	r.GET("/statefulset/list/:namespace", apis.GetStatefulSetList)
	r.GET("/daemonset/list/:namespace", apis.GetDaemonSetList)
	r.GET("/service/list/:namespace", apis.GetServiceList)
	r.GET("/configmap/list/:namespace", apis.GetConfigMapList)
	r.GET("/secret/list/:namespace", apis.GetSecretList)
	r.GET("/cronjob/list/:namespace", apis.GetCronJobList)
	r.GET("/pv/list", apis.GetPvList)
}
