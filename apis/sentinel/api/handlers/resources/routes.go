package resources

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/resources/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up resource specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/resources")

	router.GET("/", get(service))
	router.GET("/:id", injection.VerifyResourceOwnership, getByID(service))
	router.GET("/:id/policies", injection.VerifyResourceOwnership, getAllAssociatedPolicies(service))

	router.POST("/", injection.ValidateNewResource, post(service))
	router.POST("/:id/policies", injection.VerifyResourceOwnership, associatePolicy(service))

	router.PATCH("/:id", injection.VerifyResourceOwnership, patch(service))

	router.DELETE("/:id", injection.VerifyResourceOwnership, delete(service))
}
