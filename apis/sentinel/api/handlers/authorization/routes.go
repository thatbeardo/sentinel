package authorizations

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up context specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/principal")
	router.GET("/:principal_id/authorization", injection.VerifyResourceOwnership(service, "principal_id"), getAuthorizationForPrincipal(service))
}
