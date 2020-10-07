package contexts

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/service"
	"github.com/gin-gonic/gin"
)

// @Summary Update a context by it's ID
// @Tags Contexts
// @Description Update Polciy name.
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param id path string true "context ID"
// @Param input body context.Input true "New name to be assigned to an existing context"
// @Success 204 {object} context.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/contexts/{id} [patch]
func patch(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input context.Input
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		output, err := service.Update(id, &input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, output)
	}
}
