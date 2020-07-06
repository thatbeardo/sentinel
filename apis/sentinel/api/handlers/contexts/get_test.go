package contexts_test

import (
	"net/http"
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestInvalidPath(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/invalid-path/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("", "query-parameter-todo", "Path not found", http.StatusNotFound), http.StatusNotFound)
}

func TestGetByID_ServiceReturnsContexts_Return200(t *testing.T) {
	mockService := mockService{
		GetByIDResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/contexts/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.OutputDetails, http.StatusOK)
}

func TestGetByID_ServiceReturnsDatabaseError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/contexts/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/contexts/:id", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}
