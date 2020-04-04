package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "resource",
		})
	}
}
