package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/middleware/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/statsd"
	"github.com/gin-gonic/gin"
)

// VerifyToken function verifies the incoming jwt token
func VerifyToken(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-sentinel-tenant")
	if err := injection.VerifyAccessToken(c.Writer, c.Request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			views.GenerateErrorResponse(
				http.StatusUnauthorized,
				"The access token is invalid. Please provide a valid token in the header",
				c.Request.URL.Path,
			),
		)
	}
}

// VerifyTenant ensures that the user has added tenant information in the header
func VerifyTenant(c *gin.Context) {

	scope := injection.ExtractClaim(c.Request.Context(), "scope")
	if len(scope) > 0 {
		return
	}

	tenant := c.Request.Header.Get("x-sentinel-tenant")
	if tenant == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			views.GenerateErrorResponse(
				http.StatusBadRequest,
				"Please add tenant in the header.",
				c.Request.URL.Path,
			),
		)
		return
	}
	c.Request = c.Request.WithContext(injection.AddTenantToContext(c.Request.Context(), tenant))
}

// VerifyContextOwnership checks if the context being updated belongs to the correct tenant
func VerifyContextOwnership(s service.Service, identifier string) gin.HandlerFunc {
	return func(c *gin.Context) {

		scope := injection.ExtractClaim(c.Request.Context(), "scope")
		if strings.Contains(scope, "create:permission") {
			return
		}

		isValidcontext := s.IsContextOwnedByClient(c.Request.Context(), c.Param(identifier))
		if !isValidcontext {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				views.GenerateErrorResponse(
					http.StatusNotFound,
					"The requested context does not exist",
					c.Request.URL.Path,
				),
			)
		}
	}
}

// VerifyResourceOwnership checks if the caller has access to the requested resource
func VerifyResourceOwnership(s service.Service, identifier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scope := injection.ExtractClaim(c.Request.Context(), "scope")
		if strings.Contains(scope, "create:context") {
			return
		}
		validOwnership := s.IsTargetOwnedByClient(c.Request.Context(), c.Param(identifier))

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
		scope := injection.ExtractClaim(c.Request.Context(), "scope")
		if strings.Contains(scope, "create:resource") {
			return
		}

		var input resource.Input
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		err := injection.Unmarshal(bodyBytes, &input)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				views.GenerateErrorResponse(
					http.StatusBadRequest,
					"Malformed request body",
					c.Request.URL.Path,
				),
			)
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if input.Data.Relationships != nil && input.Data.Relationships.Parent != nil {
			validParent := s.IsTargetOwnedByClient(c.Request.Context(), input.Data.Relationships.Parent.Data.ID)
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

// VerifyRelationshipOwnership makes sure that the edge being created/updated belongs to the correct tenant
func VerifyRelationshipOwnership(s service.Service, identifier string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isValidPermission := s.IsPermissionOwnedByTenant(c.Request.Context(), c.Param(identifier))

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

// Metrics is a middleware responsible to report metrics to statsd
func Metrics(client statsd.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		// send call count
		client.Increment(fmt.Sprintf("%s.%s", constructPathMetric(c.Request.URL.Path), Count))

		// send response code
		status := c.Writer.Status()
		client.Increment(fmt.Sprintf("%s.%s", constructPathMetric(c.Request.URL.Path), strconv.Itoa(status)))

		// send response time
		duration := time.Since(startTime).Seconds() * 1000 // in milliseconds
		client.Timing(constructPathMetric(c.Request.URL.Path), duration)

		// send response size
		size := c.Writer.Size()
		client.Gauge(constructPathMetric(c.Request.URL.Path), size)
	}
}

func constructPathMetric(path string) string {
	return fmt.Sprintf("%s.%s.%s.%s", Organization, Class, os.Getenv("ENV"), path)
}
