package policies

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	authorizationService "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service, authorizationService authorizationService.Service) {
	router := r.Group("/policies")

	const identifier = "id"

	router.GET("/:id", injection.VerifyPolicyOwnership(authorizationService, identifier), getByID(service))

	router.PATCH("/:id", injection.VerifyPolicyOwnership(authorizationService, identifier), patch(service))

	router.DELETE("/:id", injection.VerifyPolicyOwnership(authorizationService, identifier), delete(service))

}
