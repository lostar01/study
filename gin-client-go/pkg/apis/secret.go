package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetSecretList(c *gin.Context) {
	namespace := c.Param("namespace")
	secretList, err := service.GetSecretList(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, secretList)
}
