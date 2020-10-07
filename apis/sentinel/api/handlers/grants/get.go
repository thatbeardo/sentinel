package grants

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/service"
	"github.com/gin-gonic/gin"
)

// @Summary Shows all Principal access to Target Resources managed through a context
// @Tags Grants
// @Description Shows all Principal access to Target Resources managed through a context
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param resource_id path string true "Resource ID"
// @Success 200 {object} grant.Output	"ok"
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/grants/resources/{resource_id} [get]
func getPrincipalsAndContextsForResource(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("resource_id")
		grants, err := service.GetPrincipalAndcontextForResource(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, grants)
	}
}
