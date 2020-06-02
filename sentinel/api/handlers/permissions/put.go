package permissions

import (
	"net/http"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"

	// Used for swaggo annotation
	_ "github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
)

// @Summary Update a permission that permits apolicy on a resource
// @Tags Permissions
// @Description Create a permission for a policy on a resource
// @Accept  json
// @Produce  json
// @Param policy_id path string true "Policy ID"
// @Param resource_id path string true "Resource ID"
// @Param input body inputs.Payload true "Details about the permission to be added"
// @Success 204 {object} outputs.Response	"ok"
// @Success 404 {object} views.ErrView
// @Router /v1/policies/{policy_id}/resources/{resource_id}/permissions [put]
func put(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		policyID := c.Param("policy_id")
		resourceID := c.Param("resource_id")
		var payload inputs.Payload
		if err := c.ShouldBind(&payload); err != nil {
			views.Wrap(err, c)
			return
		}
		resourceResponse, err := service.Create(&payload, policyID, resourceID)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, resourceResponse)
	}
}
