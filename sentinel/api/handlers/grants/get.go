package grants

import (
	"net/http"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/service"
	"github.com/gin-gonic/gin"
)

// @Summary Shows all Principal access to Target Resources managed through a Policy
// @Tags Grants
// @Description Shows all Principal access to Target Resources managed through a Policy
// @Accept  json
// @Produce  json
// @Param resource_id path string true "Resource ID"
// @Success 200 {object} grant.Output	"ok"
// @Success 500 {object} views.ErrView
// @Router /v1/grants/resources/{resource_id} [get]
func getPrincipalsAndPoliciesForResource(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("resource_id")
		grants, err := service.GetPrincipalAndPolicyForResource(id)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, grants)
	}
}
