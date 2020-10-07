package contexts

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/service"
	"github.com/gin-gonic/gin"
)

// @Summary Get context by ID
// @Tags Contexts
// @Description Get a context by its ID
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param id path string true "context ID"
// @Success 200 {object} context.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/contexts/{id} [get]
func getByID(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		context, err := service.GetByID(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, context)
	}
}
