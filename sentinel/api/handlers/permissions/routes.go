package permissions

import (
	"github.com/bithippie/guard-my-app/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/policies")
	router.PUT("/:policy_id/resources/:resource_id/permissions", put(service))
}
