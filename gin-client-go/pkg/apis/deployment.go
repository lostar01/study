package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func ListDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	deployments, err := service.ListDeployment(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, deployments)
}

func GetDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	deployment, err := service.GetDeployment(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, deployment)
}
