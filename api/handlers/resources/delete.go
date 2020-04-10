package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/api/views"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// @Summary Delete a resource by its ID
// @Tags Resources
// @Description Delete a resource by its ID
// @Accept  json
// @Produce  json
// @Param id path string true "Resource ID"
// @Success 204 {object} resource.Response	"ok"
// @Success 500 {object} views.ErrView
// @Router /v1/resources/{id} [delete]
func delete(service resource.Service) gin.HandlerFunc {
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
