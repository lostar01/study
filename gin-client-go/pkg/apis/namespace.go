package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lostar.com/m/pkg/service"
)

func GetNameSpace(c *gin.Context) {
	namespaces, err := service.GetNameSpace()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, namespaces)
}
