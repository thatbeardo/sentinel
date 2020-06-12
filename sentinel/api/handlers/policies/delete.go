package policies

import (
	"net/http"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// @Summary Delete a policy by its ID
// @Tags Policies
// @Description Delete a policy by its ID
// @Accept  json
// @Produce  json
// @Param id path string true "Policy ID"
// @Success 404 {object} views.ErrView
// @Router /v1/policies/{id} [delete]
func delete(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := service.Delete(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
