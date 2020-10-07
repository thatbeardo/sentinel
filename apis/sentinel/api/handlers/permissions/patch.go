package permissions

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
)

// @Summary Update a Permission by its ID
// @Tags Permissions
// @Description Update Permission Details.
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param id path string true "Permission ID"
// @Param input body permission.Input true "New name to be assigned to an existing permission"
// @Success 202 {object} permission.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/permissions/{id} [patch]
func patch(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("permission_id")
		var input permission.Input
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
