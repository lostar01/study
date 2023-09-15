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
	}
	c.JSON(http.StatusOK, deployments)
}
