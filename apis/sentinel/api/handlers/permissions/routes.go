package permissions

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	authorizationService "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// Routes sets up context specific routes on the engine instance
// Unexpected routes thanks to - https://github.com/gin-gonic/gin/issues/1681
func Routes(r *gin.RouterGroup, service service.Service, authorizationService authorizationService.Service) {
	router := r.Group("/permissions")
	router.PUT("/:context_id/resources/:resource_id",
		injection.VerifyPermissionIdempotence(service, "context_id", "resource_id"),
		injection.VerifyContextOwnership(authorizationService, "context_id"),
		injection.VerifyResourceOwnership(authorizationService, "resource_id"),
		put(service))

	router.GET("/:context_id/resources",
		injection.VerifyContextOwnership(authorizationService, "context_id"),
		getAllPermissionsForAcontext(service))

	router.GET("/:context_id/resources/:resource_id",
		injection.VerifyContextOwnership(authorizationService, "context_id"),
		injection.VerifyResourceOwnership(authorizationService, "resource_id"),
		getAllPermissionsForAcontextWithResource(service))

	router.PATCH("/:permission_id", injection.VerifyPermissionOwnership(authorizationService, "permission_id"), patch(service))
	router.DELETE("/:permission_id", injection.VerifyPermissionOwnership(authorizationService, "permission_id"), delete(service))
}
