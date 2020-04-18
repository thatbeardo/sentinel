package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// ResourceRoutes sets up resource specific routes on the engine instance
func ResourceRoutes(r *gin.RouterGroup, service resource.Service) {
	router := r.Group("/resources")

	router.GET("/", get(service))
	router.GET("/:id", getByID(service))
	router.POST("/", post(service))
	router.PATCH("/:id", patch(service))
	router.DELETE("/:id", delete(service))
}
