package grants

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	authorizationService "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service, authorizationService authorizationService.Service) {
	router := r.Group("/grants")

	router.PUT("/resources/:resource_id/policies/:policy_id",
		injection.VerifyPolicyOwnership(authorizationService, "policy_id"),
		injection.VerifyResourceOwnership(authorizationService, "resource_id"),
		put(service))

	router.GET("/resources/:resource_id",
		injection.VerifyResourceOwnership(authorizationService, "resource_id"),
		getPrincipalsAndPoliciesForResource(service))
}
