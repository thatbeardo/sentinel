package resources

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	contextDto "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-gonic/gin"
)

type key string

const (
	tenant key = "tenant"
)

// @Summary Create a new Resource
// @Description Add a new resource to existing resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired tenant - environment"
// @Param claimant query string true "Claimant requesting the operation"
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

// @Summary Associate a new context with an existing resource
// @Description Grant a new context to an existing principal resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired tenant - environment"
// @Param id path string true "Principal Resource ID"
// @Param claimant query string true "Claimant requesting the operation"
// @Param input body context.Input true "context to be created"
// @Success 202 {object} context.OutputDetails	"ok"
// @Failure 500 {object} views.ErrView	"ok"
// @Security ApiKeyAuth
// @Router /v1/resources/{id}/contexts [post]
func associatecontext(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input contextDto.Input
		principalID := c.Param("id")
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		response, err := service.Associatecontext(principalID, &input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, response)
	}
}
