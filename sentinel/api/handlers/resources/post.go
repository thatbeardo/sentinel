package resources

import (
	"net/http"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	resource "github.com/bithippie/guard-my-app/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/service"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new Resource
// @Description Add a new resource to existing resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param input body resource.Input true "Resource to be created"
// @Success 202 {object} resource.OutputDetails	"ok"
// @Failure 500 {object} views.ErrView	"ok"
// @Router /v1/resources [post]
func post(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input resource.Input
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		response, err := service.Create(&input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, response)
	}
}
