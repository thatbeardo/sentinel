package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/pkg/resource"
)

func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		resource := &resource.Resource{Name: "Harshil"}
		c.JSON(http.StatusOK, resource)
	}
}
