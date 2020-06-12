package policies

import (
	"net/http"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// @Summary Get all the Policies
// @Tags Policies
// @Description Get all the Policies stored
// @Accept  json
// @Produce  json
// @Success 200 {object} policy.Output	"ok"
// @Success 500 {object} views.ErrView
// @Router /v1/policies [get]
func get(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		policies, err := service.Get()
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, policies)
	}
}

// @Summary Get Policy by ID
// @Tags Policies
// @Description Get a Policy by its ID
// @Accept  json
// @Produce  json
// @Param id path string true "Policy ID"
// @Success 200 {object} policy.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Success 500 {object} views.ErrView
// @Router /v1/policies/{id} [get]
func getByID(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		policy, err := service.GetByID(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, policy)
	}
}
