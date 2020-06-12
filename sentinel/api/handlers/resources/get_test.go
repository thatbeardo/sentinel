package resources_test

import (
	"net/http"
	"testing"

	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/testdata"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
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
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGetResourcesDatabaseError(t *testing.T) {
	mockService := mockService{
		GetErr: models.ErrDatabase,
	}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/resources/", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func TestGetResourceByIDOk(t *testing.T) {

	mockService := mockService{
		GetByIDResponse: testdata.ModificationResult,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.ModificationResult, http.StatusOK)
}

func TestGetResourceByIDNoResourceFound(t *testing.T) {

	mockService := mockService{
		GetByIDErr: models.ErrNotFound,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/resources/:id", "query-parameter-todo", "Data not found", http.StatusNotFound), http.StatusNotFound)
}
