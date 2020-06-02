package grants

import (
	"github.com/bithippie/guard-my-app/sentinel/models/grant/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/resources")
	router.PUT("/:resource_id/grants/:policy_id", put(service))
}
