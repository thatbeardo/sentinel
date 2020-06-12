package authorizations

import (
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/principal")
	router.GET("/:principal_id/authorization", getAuthorizationForPrincipal(service))
}
