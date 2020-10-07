package middleware_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/middleware"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/middleware/injection"
	mock "github.com/bithippie/guard-my-app/apis/sentinel/mocks"
	mocks "github.com/bithippie/guard-my-app/apis/sentinel/mocks/authorization"
	mockGrants "github.com/bithippie/guard-my-app/apis/sentinel/mocks/grants"
	mockPermissions "github.com/bithippie/guard-my-app/apis/sentinel/mocks/permissions"
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

func TestInternalMetricLogging_UDPCallsAreFired_ReturnsVoid(t *testing.T) {
	const path = "/test"
	router := setupRouter()

	router.GET(path, middleware.Metrics(mock.Statsd{
		ExpectedTimingMetric: fmt.Sprintf("%s.%s.%s.%s", middleware.Organization, middleware.Class, os.Getenv("ENV"), path),
		ExpectedGaugeMetric:  fmt.Sprintf("%s.%s.%s.%s", middleware.Organization, middleware.Class, os.Getenv("ENV"), path),
		T:                    t,
	}))

	response, cleanup := testutil.PerformRequest(router, "GET", path, "")
	assert.Equal(t, http.StatusOK, response.StatusCode)
	defer cleanup()
}

func TestGrantExists_ServiceReturnsTrue_ReturnBadRequest(t *testing.T) {
	// This is not how you set path params.
	// Refer to https://stackoverflow.com/questions/43826412/how-to-properly-set-path-params-in-url-using-golang-http-client
	// Update testutil.PerformRequest()
	const path = "/test"
	const params = "?context_id=test-context-id&resource_id=test-principal-id"
	router := setupRouter()

	router.GET(path, middleware.VerifyGrantExistence(mockGrants.GrantService{
		GrantExistsResponse: true,
		ExpectedContextID:   "test-context-id",
		ExpectedPrincipalID: "test-principal-id",
		T:                   t},
		"context_id",
		"resource_id"))

	response, cleanup := testutil.PerformRequest(router, "GET", path+params, "")
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	defer cleanup()
}

func TestGrantExists_ServiceReturnsFalse_ReturnStatusOK(t *testing.T) {
	// This is not how you set path params.
	// Refer to https://stackoverflow.com/questions/43826412/how-to-properly-set-path-params-in-url-using-golang-http-client
	// Update testutil.PerformRequest()
	const path = "/test"
	const params = "?context_id=test-context-id&resource_id=test-principal-id"
	router := setupRouter()

	router.GET(path, middleware.VerifyGrantExistence(mockGrants.GrantService{
		GrantExistsResponse: false,
		ExpectedContextID:   "test-context-id",
		ExpectedPrincipalID: "test-principal-id",
		T:                   t},
		"context_id",
		"resource_id"))

	response, cleanup := testutil.PerformRequest(router, "GET", path+params, "")
	assert.Equal(t, http.StatusOK, response.StatusCode)
	defer cleanup()
}

func TestPermissionIdempotence_ErrorInUnmarshallingBody_ReturnStatusBadRequest(t *testing.T) {
	defer injection.Reset()
	injection.Unmarshal = func(data []byte, v interface{}) error {
		return errors.New("Some test error")
	}

	const path = "/test"
	router := setupRouter()
	router.POST(path, middleware.VerifyPermissionIdempotence(mockPermissions.PermissionService{
		IsPermissionIdempotentResponse: false,
	},
		"context_id",
		"resource_id"))
	response, cleanup := testutil.PerformRequest(router, "POST", path, `{"data"}`)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		testutil.GenerateError("/test", "query-parameter-todo", "Malformed request body", 400),
		http.StatusBadRequest)

}

func TestPermissionIdempotence_PermissionAlreadyExists_ReturnStatusBadRequest(t *testing.T) {
	defer injection.Reset()
	injection.Unmarshal = func(data []byte, v interface{}) error {
		return nil
	}

	const path = "/test"
	router := setupRouter()
	router.POST(path, middleware.VerifyPermissionIdempotence(mockPermissions.PermissionService{
		IsPermissionIdempotentResponse: false,
	},
		"context_id",
		"resource_id"))
	response, cleanup := testutil.PerformRequest(router, "POST", path, `{"data"}`)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		testutil.GenerateError(
			"/test",
			"query-parameter-todo",
			"A permission with the same name already exists between the context-resource pair",
			http.StatusBadRequest),
		http.StatusBadRequest)

}

func TestPermissionIdempotence_NewPermissionCreated_ReturnStatusOK(t *testing.T) {
	defer injection.Reset()
	injection.Unmarshal = func(data []byte, v interface{}) error {
		return nil
	}

	const path = "/test"
	router := setupRouter()
	router.POST(path, middleware.VerifyPermissionIdempotence(mockPermissions.PermissionService{
		IsPermissionIdempotentResponse: true,
	},
		"context_id",
		"resource_id"))
	response, cleanup := testutil.PerformRequest(router, "POST", path, `{"data"}`)
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestVerifyClaimant_ClaimantNotProvided_ReturnsStatusBadRequest(t *testing.T) {
	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyClaimant)
	response, cleanup := testutil.PerformRequest(router, "GET", path, ``)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		testutil.GenerateError(
			"/test",
			"query-parameter-todo",
			"Please provide a claimant",
			http.StatusBadRequest),
		http.StatusBadRequest)
}

func TestVerifyClaimant_ClaimantProvided_ReturnsStatusOK(t *testing.T) {
	const path = "/test"
	router := setupRouter()
	router.GET(path, middleware.VerifyClaimant)
	response, cleanup := testutil.PerformRequest(router, "GET", path+"?claimant=test", ``)
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	return r
}
