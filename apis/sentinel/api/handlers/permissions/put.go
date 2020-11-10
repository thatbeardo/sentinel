package permissions

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// @Summary Update a permission that permits acontext on a resource
// @Tags Permissions
// @Description Create a permission for a context on a resource
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param context_id path string true "context ID"
// @Param resource_id path string true "Resource ID"
// @Param input body permission.Input true "Details about the permission to be added"
// @Success 202 {object} permission.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/permissions/{context_id}/resources/{resource_id} [put]
func put(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextID := c.Param("context_id")
		resourceID := c.Param("resource_id")
		var input permission.Input
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		resourceResponse, err := service.Create(&input, contextID, resourceID)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, resourceResponse)
	}
}
