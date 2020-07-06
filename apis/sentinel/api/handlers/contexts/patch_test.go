package contexts_test

import (
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestPatch_ServiceUpdatescontext_ReturnsStatusAccepted(t *testing.T) {
	mockService := mockService{
		UpdateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/contexts/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.OutputDetails, http.StatusAccepted)
}

func TestPatch_ErrorFromService_ReturnsInternalServerError(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/contexts/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/contexts/:id", "query-parameter-todo", "Database Error", http.StatusInternalServerError), http.StatusInternalServerError)
}

func TestPatch_NameAttributeAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/contexts/test-id", nameAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/contexts/:id"), http.StatusBadRequest)
}

func TestPatch_TypeFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/contexts/test-id", typeAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/contexts/:id"), http.StatusBadRequest)
}
