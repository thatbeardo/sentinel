package permissions_test

import (
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestPatch_ServiceUpdatesPermission_ReturnsStatusAccepted(t *testing.T) {
	mockService := mockService{
		UpdateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/permissions/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.OutputDetails, http.StatusAccepted)
}

func TestPatch_ErrorFromService_ReturnsInternalServerError(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/permissions/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/permissions/:permission_id", "query-parameter-todo", "Database Error", http.StatusInternalServerError), http.StatusInternalServerError)
}

func TestPatch_NameAttributeAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/permissions/test-id", nameFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/permissions/:permission_id"), http.StatusBadRequest)
}

func TestPatch_TypeFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/permissions/test-id", typeFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/permissions/:permission_id"), http.StatusBadRequest)
}

func TestPatch_PermittedFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/permissions/test-id", permittedFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.Permitted' Error:Field validation for 'Permitted' failed on the 'required' tag", "/v1/permissions/:permission_id"), http.StatusBadRequest)
}
