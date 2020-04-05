package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/pkg/resource"
)

// ResourceRoutes sets up resource specific routes on the engine instance
func ResourceRoutes(r *gin.Engine, service resource.Service) {
	router := r.Group("/resources")
	router.GET("/", get(service))
}
