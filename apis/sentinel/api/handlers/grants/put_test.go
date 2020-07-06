package grants_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestPut_AllParametersPresent_Returns200(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/grants/resources/test-resource-id/contexts/test-context-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.OutputDetails, http.StatusAccepted)
}

func TestPut_TypeFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/grants/resources/test-resource-id/contexts/test-context-id", typeFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/grants/resources/:resource_id/contexts/:context_id"),
		http.StatusBadRequest)
}

func TestPut_WithGrantFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/grants/resources/test-resource-id/contexts/test-context-id", withGrantFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Attributes.WithGrant' Error:Field validation for 'WithGrant' failed on the 'required' tag", "/v1/grants/resources/:resource_id/contexts/:context_id"),
		http.StatusBadRequest)
}

func TestPut_ServiceReturnsError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: errors.New("some-test-error"),
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/grants/resources/test-resource-id/contexts/test-context-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "some-test-error", "/v1/grants/resources/:resource_id/contexts/:context_id"),
		http.StatusInternalServerError)
}
