package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/claims"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/gin-gonic/gin"
)

// VerifyResourceOwnership function checks if the caller has access to the requested resource
func VerifyResourceOwnership(c *gin.Context) {
	azp := claims.Extract(c.Request.Context(), "azp")
	validOwnership := isTargetOwnedByTenant(c.Param("id"), azp)

	if !validOwnership {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			views.GenerateErrorResponse(
				http.StatusNotFound,
				"The requested resource does not exist",
				c.Request.URL.Path,
			),
		)
	}
}

// ValidateNewResource checks if the scope is set, or if the parent is reachable
func ValidateNewResource(c *gin.Context) {
	scope := claims.Extract(c.Request.Context(), "scope")
	azp := claims.Extract(c.Request.Context(), "azp")
	var input resource.Input

	if strings.Contains(scope, "create:resource") {
		return
	}

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(bodyBytes, &input)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if input.Data.Relationships != nil && input.Data.Relationships.Parent != nil {
		validParent := isTargetOwnedByTenant(input.Data.Relationships.Parent.Data.ID, azp)
		if !validParent {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				views.GenerateErrorResponse(
					http.StatusNotFound,
					"The parent resource does not exist",
					c.Request.URL.Path,
				),
			)
		}
	}
}

func isTargetOwnedByTenant(target string, tenant string) bool {
	authorizationService := service.NewWithoutRepository()
	return authorizationService.IsTargetOwnedByTenant(target, tenant)
}
