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
	}
	c.JSON(http.StatusOK, podList)
}

func GetPod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	pod, err := service.GetPod(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, pod)
}
