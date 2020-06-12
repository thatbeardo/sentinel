package policies

import (
	"github.com/bithippie/guard-my-app/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/policies")
	router.GET("/", get(service))
	router.GET("/:id", getByID(service))
	router.POST("/", post(service))
	router.PATCH("/:id", patch(service))
	router.DELETE("/:id", delete(service))
}
