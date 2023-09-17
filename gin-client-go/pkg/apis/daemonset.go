package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetDaemonSetList(c *gin.Context) {
	namespace := c.Param("namespace")

	daemonSetList, err := service.GetDaemonSetList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, daemonSetList)
}
