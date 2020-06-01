package policies

import (
	"net/http"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/service"
	"github.com/gin-gonic/gin"

	// Used in annotations
	_ "github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
)

// @Summary Create a new Policy
// @Description Add a new Policy to existing Policies
// @Tags Policies
// @Accept  json
// @Produce  json
// @Param input body inputs.Payload true "Policy to be created"
// @Success 202 {object} outputs.Policy	"ok"
// @Failure 500 {object} views.ErrView	"ok"
// @Router /v1/policies [post]
func post(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resourceInput inputs.Payload
		if err := c.ShouldBind(&resourceInput); err != nil {
			views.Wrap(err, c)
			return
		}
		resourceResponse, err := service.Create(&resourceInput)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusAccepted, resourceResponse)
	}
}
