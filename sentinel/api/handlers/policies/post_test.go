package policies_test

import (
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	errors "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/testdata"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
)

func TestPost_AllParametersPresent_Returns200(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.OutputDetails, http.StatusAccepted)
}

func TestPost_NameAttributeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", nameAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
}

func TestPost_AttributeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", nameAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
}

func TestPost_TypeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", typeAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
}

func TestPost_DataAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", dataAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Input.Data' Error:Field validation for 'Data' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
}

func TestPost_ServiceReturnsError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: errors.ErrDatabase,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Database Error", "/v1/policies/"), http.StatusInternalServerError)
}
