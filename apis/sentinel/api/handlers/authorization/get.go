package authorizations

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/gin-gonic/gin"
)

// @Summary Get authorization details about a principal.
// @Tags Authorization
// @Description Shows all the permissions this principal has to provided target resources
// @Accept  json
// @Produce  json
// @Param x-sentinel-tenant header string true "Desired environment"
// @Param claimant query string true "Claimant requesting the operation"
// @Param principal_id path string true "Principal ID"
// @Param permissions query []string false "Name of the permissions which allow access to the target"
// @Param targets query []string false "Name of the targtes to which a permission allows access"
// @Param context_id query string false "Context through which authorization is determined"
// @Param depth query int false "Limit your search results." default(0)
// @Param include_denied query boolean false "Include permissions that have deny permit fields set" default(false)
// @Success 200 {object} authorization.Output	"ok"
// @Success 500 {object} views.ErrView
// @Security ApiKeyAuth
// @Router /v1/principal/{principal_id}/authorization [get]
func getAuthorizationForPrincipal(service service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		principalID := c.Param("principal_id")
		depth, err := strconv.Atoi(c.Query("depth"))

		var targets []string
		if len(c.Query("targets")) > 0 {
			targets = strings.Split(c.Query("targets"), ",")
		}

		var permissions []string
		if len(c.Query("permissions")) > 0 {
			permissions = strings.Split(c.Query("permissions"), ",")
		}

		if err != nil {
			depth = 0
		}

		includeDenied, err := strconv.ParseBool(c.Query("include_denied"))
		if err != nil {
			includeDenied = false
		}

		var input = authorization.Input{
			Targets:       targets,
			Permissions:   permissions,
			Depth:         depth,
			IncludeDenied: includeDenied,
		}

		contextID := c.Query("context_id")
		authorization, err := service.GetAuthorizationForPrincipal(principalID, contextID, input)
		if err != nil {
			views.Wrap(err, c)
			return
		}
		c.JSON(http.StatusOK, authorization)
	}
}
