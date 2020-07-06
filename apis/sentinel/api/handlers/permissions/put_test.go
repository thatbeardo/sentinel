package permissions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestPut_AllParametersPresent_Returns200(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/permissions/test-context-id/resources/test-target-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.OutputDetails, http.StatusAccepted)
}

func TestPut_TypeFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/permissions/test-context-id/resources/test-target-id", typeFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/permissions/:context_id/resources/:resource_id"),
		http.StatusBadRequest)
}

func TestPut_NameFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/permissions/test-context-id/resources/test-target-id", nameFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/permissions/:context_id/resources/:resource_id"),
		http.StatusBadRequest)
}

func TestPut_PermittedFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/permissions/test-context-id/resources/test-target-id", permittedFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Attributes.Permitted' Error:Field validation for 'Permitted' failed on the 'required' tag", "/v1/permissions/:context_id/resources/:resource_id"),
		http.StatusBadRequest)
}

func TestPut_ServiceReturnsError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: errors.New("some-test-error"),
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/permissions/test-context-id/resources/test-target-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "some-test-error", "/v1/permissions/:context_id/resources/:resource_id"),
		http.StatusInternalServerError)
}
