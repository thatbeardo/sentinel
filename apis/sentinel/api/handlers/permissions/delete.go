package permissions

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// @Summary Delete a permission by its ID
// @Tags Permissions
// @Description Delete a permission by its ID
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param id path string true "Permission ID"
// @Success 404 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/permissions/{id} [delete]
func delete(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("permission_id")
		err := service.Delete(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
