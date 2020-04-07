package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// @Summary Get all the resources
// @Description Get all the nodes present in the graph
// @Accept  json
// @Produce  json
// @Success 200 {object} resource.Response	"ok"
// @Router /resources [get]
func get(service resource.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resources, _ := service.Get()
		c.JSON(http.StatusOK, resources)
	}
}
