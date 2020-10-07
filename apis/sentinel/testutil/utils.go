package testutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// PerformRequestWithQueryParams creates a request with query parameters
func PerformRequestWithQueryParams(r http.Handler, method, path, body string, queryParams map[string]string) (*http.Response, func() error) {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	params := req.URL.Query()
	for k, v := range queryParams {
		params.Add(k, v)
	}
	req.URL.RawQuery = params.Encode()
	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)
	response := recorder.Result()
	return response, response.Body.Close
}

// PerformRequest creates and returns an initialized ResponseRecorder
func PerformRequest(r http.Handler, method, path, body string) (*http.Response, func() error) {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)
	response := recorder.Result()
	return response, response.Body.Close
}

// ValidateResponse reads and asserts the response body content
func ValidateResponse(t *testing.T, response *http.Response, expected interface{}, code int) {
	assert.Equal(t, code, response.StatusCode)
	b, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	out, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(out), string(b))
}

// GenerateError creates an error response instance
func GenerateError(pointer string, parameter string, detail string, code int) views.ErrView {
	source := views.Source{
		Pointer:   pointer,
		Parameter: parameter,
	}
	return views.ErrView{
		ID:     "error-id-todo",
		Status: code,
		Source: source,
		Detail: detail,
	}
}

var noMiddleware = func(c *gin.Context) {}

// RemoveMiddleware injects middlewares that don't carry out any verification
func RemoveMiddleware() {
	injection.VerifyResourceOwnership = func(authorization.Service, string) gin.HandlerFunc {
		return noMiddleware
	}

	injection.VerifyContextOwnership = func(authorization.Service, string) gin.HandlerFunc {
		return noMiddleware
	}

	injection.ValidateNewResource = func(authorization.Service) gin.HandlerFunc {
		return noMiddleware
	}

	injection.VerifyPermissionOwnership = func(authorization.Service, string) gin.HandlerFunc {
		return noMiddleware
	}

	injection.VerifyPermissionIdempotence = func(permission.Service, string, string) gin.HandlerFunc {
		return noMiddleware
	}
}
