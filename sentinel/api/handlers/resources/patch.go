package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/bithippie/go-sentinel/sentinel/api/views"
	entity "github.com/bithippie/go-sentinel/sentinel/models/resource"
	"github.com/bithippie/go-sentinel/sentinel/models/resource/service"
)

// @Summary Update a resource by it's ID
// @Tags Resources
// @Description Update resource name, sourceID, parent, etc
// @Accept  json
// @Produce  json
// @Param id path string true "Resource ID"
// @Param input body entity.Input true "Resource to be created"
// @Success 204 {object} entity.Response	"ok"
// @Success 404 {object} views.ErrView
// @Router /v1/resources/{id} [patch]
func patch(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var resourceInput entity.Input
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
