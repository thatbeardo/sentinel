package resources

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-gonic/gin"
)

// @Summary Get all the resources
// @Tags Resources
// @Description Get all the resources stored
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired tenant - environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Success 200 {object} resource.Output	"ok"
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/resources [get]
func get(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resources, err := service.Get(c.Request.Context())
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, resources)
	}
}

// @Summary Get resource by ID
// @Tags Resources
// @Description Get a resource by its ID
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired tenant - environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param id path string true "Resource ID"
// @Success 200 {object} resource.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/resources/{id} [get]
func getByID(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		resource, err := service.GetByID(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, resource)
	}
}

// @Summary Get all contexts granted to this resource
// @Tags Resources
// @Description Get all context and details that are granted to this context
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired tenant - environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param id path string true "Resource ID"
// @Success 200 {object} context.Output	"ok"
// @Success 404 {object} views.ErrView
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/resources/{id}/contexts [get]
func getAllAssociatedContexts(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		contexts, err := service.GetAllAssociatedContexts(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, contexts)
	}
}
