package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "Sentinel online")
	}
}
