package policies_test

import (
	"net/http"
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/testdata"
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
		GetResponse: testdata.Output,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/policies/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGetPoliciesDatabaseError(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/policies/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/policies/", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func TestGetByID_ServiceReturnsPolicies_Return200(t *testing.T) {
	mockService := mockService{
		GetByIDResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/policies/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.OutputDetails, http.StatusOK)
}

func TestGetByID_ServiceReturnsDatabaseError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/policies/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/policies/:id", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}
