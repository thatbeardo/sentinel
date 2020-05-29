package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/bithippie/guard-my-app/sentinel/api/views"
	entity "github.com/bithippie/guard-my-app/sentinel/models/resource"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/service"
)

// @Summary Create a new Resource
// @Description Add a new resource to existing resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param input body entity.Input true "Resource to be created"
// @Success 202 {object} entity.Element	"ok"
// @Failure 500 {object} views.ErrView	"ok"
// @Router /v1/resources [post]
func post(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resourceInput entity.Input
		if err := c.ShouldBind(&resourceInput); err != nil {
			views.Wrap(err, c)
			return
		}
		resourceResponse, err := service.Create(&resourceInput)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, resourceResponse)
	}
}
