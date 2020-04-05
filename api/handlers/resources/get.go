package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	views "github.com/thatbeardo/go-sentinel/api/views/responses/resources"
	"github.com/thatbeardo/go-sentinel/pkg/resource"
)

func get(service resource.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resource, _ := service.Get()
		c.JSON(http.StatusOK, views.WrapGetResource(*resource))
	}
}
