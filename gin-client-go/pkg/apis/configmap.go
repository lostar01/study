package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetConfigMapList(c *gin.Context) {
	namespace := c.Param("namespace")
	configMapList, err := service.GetConfigMapList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, configMapList)
}
