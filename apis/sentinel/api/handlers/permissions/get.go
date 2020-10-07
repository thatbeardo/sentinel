package permissions

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// @Summary List all Permissions for all Target Resources in context.
// @Tags Permissions
// @Description List all Permissions for all Target Resources in context.
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param context_id path string true "context ID"
// @Success 200 {object} permission.Output	"ok"
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/permissions/{context_id}/resources [get]
func getAllPermissionsForAcontext(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("context_id")
		permissions, err := service.GetAllPermissionsForcontext(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, permissions)
	}
}

// @Summary List all permissions for a context for a given target
// @Tags Permissions
// @Description List all permissions for a context for a given target
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param context_id path string true "context ID"
// @Param resource_id path string true "Resource ID"
// @Success 200 {object} permission.Output	"ok"
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/permissions/{context_id}/resources/{resource_id} [get]
func getAllPermissionsForAcontextWithResource(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextID := c.Param("context_id")
		resourceID := c.Param("resource_id")
		permissions, err := service.GetAllPermissionsForcontextWithResource(contextID, resourceID)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, permissions)
	}
}
