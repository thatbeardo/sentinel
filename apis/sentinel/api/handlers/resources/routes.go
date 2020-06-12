package resources

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up resource specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/resources")

	router.GET("/", get(service))
	router.GET("/:id", getByID(service))
	router.POST("/", post(service))
	router.PATCH("/:id", patch(service))
	router.DELETE("/:id", delete(service))
}
