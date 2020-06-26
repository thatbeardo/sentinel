package middleware

import (
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/injection"
	"github.com/gin-gonic/gin"
)

// VerifyRelationshipOwnership makes sure that the edge being created/updated belongs to the correct tenant
func VerifyRelationshipOwnership(s service.Service, identifier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := injection.ExtractClaims(c.Request.Context(), "azp")
		isValidPermission := s.IsPermissionOwnedByTenant(c.Param(identifier), tenantID)

		if !isValidPermission {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				views.GenerateErrorResponse(
					http.StatusNotFound,
					"The permission you are trying to update does not exist",
					c.Request.URL.Path,
				),
			)
		}
	}
}
