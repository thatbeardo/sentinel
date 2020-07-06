package contexts

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	authorizationService "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up context specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service, authorizationService authorizationService.Service) {
	router := r.Group("/contexts")

	const identifier = "id"

	router.GET("/:id", injection.VerifyContextOwnership(authorizationService, identifier), getByID(service))

	router.PATCH("/:id", injection.VerifyContextOwnership(authorizationService, identifier), patch(service))

	router.DELETE("/:id", injection.VerifyContextOwnership(authorizationService, identifier), delete(service))

}
