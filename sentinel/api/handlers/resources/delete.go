package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/service"
)

// @Summary Delete a resource by its ID
// @Tags Resources
// @Description Delete a resource by its ID
// @Accept  json
// @Produce  json
// @Param id path string true "Resource ID"
// @Success 204 {object} entity.Response	"ok"
// @Success 404 {object} views.ErrView
// @Router /v1/resources/{id} [delete]
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
