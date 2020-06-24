package policies

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// @Summary Get Policy by ID
// @Tags Policies
// @Description Get a Policy by its ID
// @Accept  json
// @Produce  json
// @Param id path string true "Policy ID"
// @Success 200 {object} policy.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
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
