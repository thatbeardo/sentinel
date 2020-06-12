package policies

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// @Summary Update a Policy by it's ID
// @Tags Policies
// @Description Update Polciy name.
// @Accept  json
// @Produce  json
// @Param id path string true "Policy ID"
// @Param input body policy.Input true "New name to be assigned to an existing policy"
// @Success 204 {object} policy.OutputDetails	"ok"
// @Success 404 {object} views.ErrView
// @Router /v1/policies/{id} [patch]
func patch(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input policy.Input
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
