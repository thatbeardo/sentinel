package resources

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	authorizationService "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up resource specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service, authorizationService authorizationService.Service) {
	router := r.Group("/resources")
	const identifier = "id"

	router.GET("/", get(service))
	router.GET("/:id", injection.VerifyResourceOwnership(authorizationService, identifier), getByID(service))
	router.GET("/:id/contexts", injection.VerifyResourceOwnership(authorizationService, identifier), getAllAssociatedContexts(service))

	router.POST("/", injection.ValidateNewResource(authorizationService), post(service))
	router.POST("/:id/contexts", injection.VerifyResourceOwnership(authorizationService, identifier), associatecontext(service))

	router.PATCH("/:id", injection.VerifyResourceOwnership(authorizationService, identifier), injection.ValidateNewResource(authorizationService), patch(service))

	router.DELETE("/:id", injection.VerifyResourceOwnership(authorizationService, identifier), delete(service))

}
