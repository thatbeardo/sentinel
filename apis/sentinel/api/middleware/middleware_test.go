package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/middleware"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/middleware/injection"
	mocks "github.com/bithippie/guard-my-app/apis/sentinel/mocks/authorization"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVerifyJWT_VerificationFails_ReturnStatusUnauthorized(t *testing.T) {
	defer injection.Reset()
	injection.VerifyAccessToken = func(w gin.ResponseWriter, r *http.Request) error {
		return errors.New("test-error")
	}
	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyToken)
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	testutil.ValidateResponse(t, response,
		testutil.GenerateError(
			"/test", "query-parameter-todo",
			"The access token is invalid. Please provide a valid token in the header",
			401),
		http.StatusUnauthorized)
}

func TestVerifyTenant_CalledByGuardMyApp_DoNotPerformValidation(t *testing.T) {
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "create:resource"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyTenant)
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyTenant_HeaderSet_SetsTenantInHeader(t *testing.T) {

	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return ""
	}
	const path = "/test"
	recorder := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	_, router := gin.CreateTestContext(recorder)
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-sentinel-tenant", "development")

	router.Use(middleware.VerifyTenant)
	router.GET(path)
	router.ServeHTTP(recorder, req)

	response := recorder.Result()
	assert.Equal(t, response.StatusCode, http.StatusOK)
}

func TestVerifyTenant_HeaderNotSet_ReturnsBadRequest(t *testing.T) {
	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyTenant)
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		testutil.GenerateError(
			"/test",
			"query-parameter-todo",
			"Please add tenant in the header.", http.StatusBadRequest), http.StatusBadRequest)
}

func TestVerifyContextOwnership_ScopeContainsCreatePermission_ReturnWithoutValidation(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "create:permission"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyContextOwnership(mocks.AuthorizationService{}, "policy_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path+"?context_id=test", "")
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyContextOwnership_AuthorizationServiceReturnsFalse_ReturnStatusNotFound(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "some-tenant-id"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyContextOwnership(mocks.AuthorizationService{}, "policy_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	assert.Equal(t, response.StatusCode, 404)
	testutil.ValidateResponse(t, response, testutil.GenerateError("/test", "query-parameter-todo", "The requested context does not exist", 404), http.StatusNotFound)
}

func TestVerifyContextOwnership_AuthorizationServiceReturnsTrue_ReturnStatusNotFound(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "some-tenant-id"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyContextOwnership(mocks.AuthorizationService{IsContextOwnedByClientResponse: true}, "policy_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyResourceOwnership_ScopeContainsCreateResource_ReturnWithoutValidation(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "create:context"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyResourceOwnership(mocks.AuthorizationService{}, "policy_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyResourceOwnership_AuthorizationServiceReturnsFalse_ReturnStatusNotFound(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "some-tenant-id"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyResourceOwnership(mocks.AuthorizationService{}, "policy_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	assert.Equal(t, response.StatusCode, 404)
	testutil.ValidateResponse(t, response, testutil.GenerateError("/test", "query-parameter-todo", "The requested resource does not exist", 404), http.StatusNotFound)
}

func TestVerifyResourceOwnership_AuthorizationServiceReturnsTrue_ReturnStatusNotFound(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "some-tenant-id"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyResourceOwnership(mocks.AuthorizationService{IsTargetOwnedByClientResponse: true}, "resource_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyNewResource_ScopeContainsCreateResource_ReturnWithoutValidation(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return "create:resource"
	}

	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.ValidateNewResource(mocks.AuthorizationService{IsTargetOwnedByClientResponse: true}))
	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyNewResource_ParentNotProvided_ReturnWithoutValidation(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return ""
	}

	const path = "/test"
	router := setupRouter()
	router.POST(path, middleware.ValidateNewResource(mocks.AuthorizationService{IsTargetOwnedByClientResponse: false}))
	response, cleanup := testutil.PerformRequest(router, "POST", path, `{"data":{"type":"resource","attributes":{"name": "test-name", "source_id":"test-id"},"relationships":{}}}`)
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyNewResource_ValidParentProvided_ReturnWithoutValidation(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return ""
	}

	const path = "/test"
	router := setupRouter()
	router.POST(path, middleware.ValidateNewResource(mocks.AuthorizationService{IsTargetOwnedByClientResponse: true}))
	response, cleanup := testutil.PerformRequest(router, "POST", path, `{"data":{"type":"resource","attributes":{"name": "test-name", "source_id":"test-id"},"relationships":{"parent":{"data":{"id": "invalid-parent-id", "type":"resource"}}}}}`)
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyNewResource_InvalidParentProvided_ReturnWithoutValidation(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return ""
	}

	const path = "/test"
	router := setupRouter()
	router.POST(path, middleware.ValidateNewResource(mocks.AuthorizationService{IsTargetOwnedByClientResponse: false}))
	response, cleanup := testutil.PerformRequest(router, "POST", path, `{"data":{"type":"resource","attributes":{"name": "test-name", "source_id":"test-id"},"relationships":{"parent":{"data":{"id": "invalid-parent-id", "type":"resource"}}}}}`)
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/test", "query-parameter-todo", "The parent resource does not exist", 404), http.StatusNotFound)
}

func TestVerifyNewResource_UnmarshallingFails_ReturnBadRequest(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return ""
	}
	injection.Unmarshal = func(data []byte, v interface{}) error {
		return errors.New("Some test error")
	}

	const path = "/test"
	router := setupRouter()
	router.POST(path, middleware.ValidateNewResource(mocks.AuthorizationService{IsTargetOwnedByClientResponse: false}))
	response, cleanup := testutil.PerformRequest(router, "POST", path, `{"data"}`)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		testutil.GenerateError("/test", "query-parameter-todo", "Malformed request body", 400),
		http.StatusBadRequest)
}

func TestVerifyRelationshipOwnership_AuthorizationServiceReturnsFalse_ReturnStatusNotFound(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return ""
	}
	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyRelationshipOwnership(mocks.AuthorizationService{IsPermissionOwnedByTenantResponse: false}, "grant_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path+"?grant_id=grant", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/test", "query-parameter-todo", "The permission you are trying to update does not exist", 404), http.StatusNotFound)
}

func TestVerifyRelationshipOwnership_AuthorizationServiceReturnsTrue_ReturnStatusOK(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaim = func(ctx context.Context, claim string) string {
		return ""
	}
	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyRelationshipOwnership(mocks.AuthorizationService{IsPermissionOwnedByTenantResponse: true}, "grant_id"))
	response, cleanup := testutil.PerformRequest(router, "GET", path+"?grant_id=grant", "")
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	return r
}
