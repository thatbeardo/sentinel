package permissions

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up policy specific routes on the engine instance
// Unexpected routes thanks to - https://github.com/gin-gonic/gin/issues/1681
func Routes(r *gin.RouterGroup, service service.Service) {
	router := r.Group("/permissions")
	router.PUT("/:policy_id/resources/:resource_id", put(service))
	router.GET("/:policy_id/resources", getAllPermissionsForAPolicy(service))
	router.GET("/:policy_id/resources/:resource_id", getAllPermissionsForAPolicyWithResource(service))
	router.PATCH("/:permission_id", patch(service))
	router.DELETE("/:permission_id", delete(service))
}
