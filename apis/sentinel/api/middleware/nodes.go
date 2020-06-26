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

// VerifyPolicyOwnership checks if the policy being updated belongs to the correct tenant
func VerifyPolicyOwnership(s service.Service, identifier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		azp := claims.Extract(c.Request.Context(), "azp")
		isValidPolicy := s.IsPolicyOwnedByTenant(c.Param(identifier), azp)

		if !isValidPolicy {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				views.GenerateErrorResponse(
					http.StatusNotFound,
					"The requested policy does not exist",
					c.Request.URL.Path,
				),
			)
		}
	}
}

// VerifyResourceOwnership checks if the caller has access to the requested resource
func VerifyResourceOwnership(s service.Service, identifier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scope := claims.Extract(c.Request.Context(), "scope")
		if strings.Contains(scope, "create:policy") {
			return
		}

		azp := claims.Extract(c.Request.Context(), "azp")
		validOwnership := s.IsTargetOwnedByTenant(c.Param(identifier), azp)

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
}

// ValidateNewResource checks if the scope is set, or if the parent is reachable
func ValidateNewResource(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		scope := claims.Extract(c.Request.Context(), "scope")
		if strings.Contains(scope, "create:resource") {
			return
		}

		azp := claims.Extract(c.Request.Context(), "azp")
		var input resource.Input

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			json.Unmarshal(bodyBytes, &input)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if input.Data.Relationships != nil && input.Data.Relationships.Parent != nil {
			validParent := s.IsTargetOwnedByTenant(input.Data.Relationships.Parent.Data.ID, azp)
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
}
