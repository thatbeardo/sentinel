package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// @Summary Create a new Resource
// @Description Add a new resource to existing resources
// @Accept  json
// @Produce  json
// @Param input body resource.Input true "Resource to be created"
// @Success 202 {object} resource.Response	"ok"
// @Router /resources [post]
func post(service resource.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceInput := &resource.Input{}
		if err := c.ShouldBind(resourceInput); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		resourceResponse, _ := service.Create(resourceInput)
		c.JSON(http.StatusAccepted, resourceResponse)
	}
}
