package permissions

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// @Summary List all Permissions for all Target Resources in Policy.
// @Tags Permissions
// @Description List all Permissions for all Target Resources in Policy.
// @Accept  json
// @Produce  json
// @Param policy_id path string true "Policy ID"
// @Success 200 {object} permission.Output	"ok"
// @Success 500 {object} views.ErrView
// @Router /v1/permissions/{policy_id}/resources [get]
func getAllPermissionsForAPolicy(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("policy_id")
		permissions, err := service.GetAllPermissionsForPolicy(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, permissions)
	}
}

// @Summary List all permissions for a policy for a given target
// @Tags Permissions
// @Description List all permissions for a policy for a given target
// @Accept  json
// @Produce  json
// @Param policy_id path string true "Policy ID"
// @Param resource_id path string true "Resource ID"
// @Success 200 {object} permission.Output	"ok"
// @Success 500 {object} views.ErrView
// @Router /v1/permissions/{policy_id}/resources/{resource_id} [get]
func getAllPermissionsForAPolicyWithResource(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		policyID := c.Param("policy_id")
		resourceID := c.Param("resource_id")
		permissions, err := service.GetAllPermissionsForPolicyWithResource(policyID, resourceID)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, permissions)
	}
}
