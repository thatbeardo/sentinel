package authorizations_test

import (
	"net/http"
	"testing"

	models "github.com/bithippie/guard-my-app/sentinel/models"
	authorization "github.com/bithippie/guard-my-app/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/testdata"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
)

func TestInvalidPath(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/invalid-path/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("", "query-parameter-todo", "Path not found", http.StatusNotFound), http.StatusNotFound)
}

func TestGetAllTargets_MultiplePermissions_ReturnsStatusOK(t *testing.T) {
	expectedInput := authorization.Input{
		Permissions: []string{"one", "two"},
		Depth:       4,
	}

	mockService := mockService{
		GetAuthorizationForPrincipalResponse: testdata.Output,
		ExpectedInput:                        expectedInput,
		t:                                    t,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/principal/principal-id/authorization?permissions=one,two&depth=4&include_denied=false", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGetAllTargets_MultipleTarget_ReturnsStatusOK(t *testing.T) {
	expectedInput := authorization.Input{
		Targets: []string{"one", "two"},
	}

	mockService := mockService{
		GetAuthorizationForPrincipalResponse: testdata.Output,
		ExpectedInput:                        expectedInput,
		t:                                    t,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/principal/principal-id/authorization?targets=one,two&depth=0&include_denied=false", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGetAllTargets_IncludeDeniedAndDepthNotProvided_ReturnsStatusOK(t *testing.T) {
	expectedInput := authorization.Input{
		Targets: []string{"one", "two"},
	}

	mockService := mockService{
		GetAuthorizationForPrincipalResponse: testdata.Output,
		ExpectedInput:                        expectedInput,
		t:                                    t,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/principal/principal-id/authorization?targets=one,two", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Output, http.StatusOK)
}

func TestGetAllTargets_ServiceReturnsError_ReturnsStatusOK(t *testing.T) {
	expectedInput := authorization.Input{
		Targets: []string{"one", "two"},
	}

	mockService := mockService{
		Err:           models.ErrDatabase,
		ExpectedInput: expectedInput,
		t:             t,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/principal/principal-id/authorization?targets=one,two", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, testutil.GenerateError("/v1/principal/:principal_id/authorization", "query-parameter-todo", "Database Error", http.StatusInternalServerError), http.StatusInternalServerError)
}
