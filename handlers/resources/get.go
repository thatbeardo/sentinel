package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/pkg/resource"
	views "github.com/thatbeardo/go-sentinel/views/responses/resources"
)

func get(service resource.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resource, _ := service.Get()
		c.JSON(http.StatusOK, views.WrapGetResource(*resource))
	}
}
