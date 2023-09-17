package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetCronJobList(c *gin.Context) {
	namespace := c.Param("namespace")
	cronJobList, err := service.GetCronJobList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, cronJobList)
}
