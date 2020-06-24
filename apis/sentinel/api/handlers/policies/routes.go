package policies

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/policies")

	router.GET("/:id", getByID(service))

	router.PATCH("/:id", patch(service))

	router.DELETE("/:id", delete(service))
}
