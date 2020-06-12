package resources_test

import (
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestPatchResourceOk(t *testing.T) {
	mockService := mockService{
		UpdateResponse: testdata.ModificationResult,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/resources/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.ModificationResult, http.StatusAccepted)
}

func TestPatchResourceDatabaseError(t *testing.T) {
	mockService := mockService{
		UpdateErr: models.ErrDatabase,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/resources/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/resources/:id", "query-parameter-todo", "Database Error", http.StatusInternalServerError), http.StatusInternalServerError)
}

func TestPatchResourcesSourceIdBlank(t *testing.T) {
	mockService := mockService{}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/resources/test-id", sourceIdBlank)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.SourceID' Error:Field validation for 'SourceID' failed on the 'required' tag", "/v1/resources/:id"), http.StatusBadRequest)
}
