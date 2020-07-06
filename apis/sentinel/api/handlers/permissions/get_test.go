package permissions_test

import (
	"net/http"
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestInvalidPath(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/invalid-path/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("", "query-parameter-todo", "Path not found", http.StatusNotFound), http.StatusNotFound)
}

func TestGetResourcesOk(t *testing.T) {
	mockService := mockService{
		GetAllPermissionsForcontextResponse: testdata.Output,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/permissions/test-id/resources", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGetContextsDatabaseError(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/permissions/test-id/resources", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/permissions/:context_id/resources", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func TestGetAllPermissions_ServiceReturnsPermissions_Return200(t *testing.T) {
	mockService := mockService{
		GetAllPermissionsForcontextResponse: testdata.Output,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/permissions/test-id/resources", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGetAllPermissions_ServiceReturnsDatabaseError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/permissions/test-id/resources", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/permissions/:context_id/resources", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func TestAllPermissionsForResource_ServiceReturnsPermissions_Return200(t *testing.T) {
	mockService := mockService{
		GetAllPermissionsForcontextWithResourceResponse: testdata.Output,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/permissions/context-id/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestAllPermissionsForResource_ServiceReturnsDatabaseError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/permissions/context-id/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/permissions/:context_id/resources/:resource_id", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}
