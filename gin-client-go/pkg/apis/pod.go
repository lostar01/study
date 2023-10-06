package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetPodList(c *gin.Context) {
	namespace := c.Param("namespace")
	podList, err := service.GetPodList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, podList)
}

func GetPod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	pod, err := service.GetPod(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, pod)
}

func ExecContainer(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.Param("containerName")
	method := c.DefaultQuery("action", "sh")
	err := service.WebSSH(namespaceName, podName, containerName, method, c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
}
