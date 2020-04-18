package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/api/views"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// @Summary Update a resource by it's ID
// @Tags Resources
// @Description Update resource name, sourceID, parent, etc
// @Accept  json
// @Produce  json
// @Param id path string true "Resource ID"
// @Param input body resource.Input true "Resource to be created"
// @Success 204 {object} resource.Response	"ok"
// @Success 404 {object} views.ErrView
// @Router /v1/resources/{id} [patch]
func patch(service resource.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var resourceInput resource.Input
		if err := c.ShouldBind(&resourceInput); err != nil {
			views.Wrap(err, c)
			return
		}
		resourceResponse, err := service.Update(id, &resourceInput)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, resourceResponse)
	}
}
