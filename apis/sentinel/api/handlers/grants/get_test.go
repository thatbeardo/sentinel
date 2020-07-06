package grants_test

import (
	"net/http"
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestInvalidPath(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/invalid-path/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("", "query-parameter-todo", "Path not found", http.StatusNotFound), http.StatusNotFound)
}

func TestGet_ServiceReturnsGrants_ReturnStatusOK(t *testing.T) {
	mockService := mockService{
		GetPrincipalAndcontextForResourceResponse: testdata.Output,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/grants/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGet_ServiceReturnsDatabaseError_ReturnInternalServerError(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/grants/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/grants/resources/:resource_id", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}
