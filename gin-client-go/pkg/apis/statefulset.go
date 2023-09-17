package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetStatefulSetList(c *gin.Context) {
	namespace := c.Param("namespace")
	statefulSetList, err := service.GetStatefulSetList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, statefulSetList)
}
