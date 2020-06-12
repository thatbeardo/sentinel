package policies

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new Policy
// @Description Add a new Policy to existing Policies
// @Tags Policies
// @Accept  json
// @Produce  json
// @Param input body policy.Input true "Policy to be created"
// @Success 202 {object} policy.OutputDetails	"ok"
// @Failure 500 {object} views.ErrView	"ok"
// @Router /v1/policies [post]
func post(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input policy.Input
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}
		response, err := service.Create(&input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, response)
	}
}
