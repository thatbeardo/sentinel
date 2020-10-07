package grants

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/service"
	"github.com/gin-gonic/gin"
)

// @Summary Update a grant that permits a context on a resource
// @Tags Grants
// @Description Create a grant for a context on a resource
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param context_id path string true "context ID"
// @Param resource_id path string true "Resource ID"
// @Param input body grant.Input true "Details about the Grant to be added"
// @Success 204 {object} grant.Output	"ok"
// @Success 404 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/grants/resources/{resource_id}/contexts/{context_id} [put]
func put(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextID := c.Param("context_id")
		resourceID := c.Param("resource_id")
		var input grant.Input
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
