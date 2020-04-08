package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/api/views"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// @Summary Get all the resources
// @Tags Resources
// @Description Get all the nodes present in the graph
// @Accept  json
// @Produce  json
// @Success 200 {object} resource.Response	"ok"
// @Success 500 {object} views.ErrView
// @Router /v1/resources [get]
func get(service resource.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resources, err := service.Get()
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, resources)
	}
}
