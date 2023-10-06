package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetServiceList(c *gin.Context) {
	namespace := c.Param("namespace")
	serviceList, err := service.GetServiceList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, serviceList)
}
