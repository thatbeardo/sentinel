package resources

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-gonic/gin"
)

// @Summary Update a resource by it's ID
// @Tags Resources
// @Description Update resource name, sourceID, parent, etc
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param id path string true "Resource ID"
// @Param claimant query string true "Claimant requesting the operation"
// @Param input body resource.Input true "Resource to be created"
// @Success 204 {object} resource.Output	"ok"
// @Success 404 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/resources/{id} [patch]
func patch(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input resource.Input
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		resource, err := service.Update(id, &input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, resource)
	}
}
