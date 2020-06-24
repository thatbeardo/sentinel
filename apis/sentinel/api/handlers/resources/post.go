package resources

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
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
// @Security ApiKeyAuth
// @Router /v1/resources/ [post]
func post(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input resource.Input
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		response, err := service.Create(c.Request.Context(), &input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, response)
	}
}

// @Summary Associate a new Policy with an existing resource
// @Description Grant a new Policy to an existing principal resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param id path string true "Principal Resource ID"
// @Param input body policy.Input true "Policy to be created"
// @Success 202 {object} policy.OutputDetails	"ok"
// @Failure 500 {object} views.ErrView	"ok"
// @Security ApiKeyAuth
// @Router /v1/resources/{id}/policies [post]
func associatePolicy(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input policy.Input
		principalID := c.Param("id")
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		response, err := service.AssociatePolicy(principalID, &input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, response)
	}
}
