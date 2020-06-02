package grants_test

import (
	"errors"
	"net/http"
	"testing"

	handler "github.com/bithippie/guard-my-app/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/sentinel/api/handlers/grants"
	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/service"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/testdata"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type mockService struct {
	CreateResponse outputs.Grant
	Err            error
}

func (m mockService) Create(*inputs.Payload, string, string) (outputs.Grant, error) {
	return m.CreateResponse, m.Err
}

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	grants.Routes(group, s)
	return r
}

const noErrors = `{"data":{"type":"grant","attributes":{"with_grant":true}}}`
const typeFieldAbsent = `{"data":{"attributes":{"with_grant":true}}}`
const withGrantFieldAbsent = `{"data":{"type":"grant","attributes":{}}}`

func TestPut_AllParametersPresent_Returns200(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/resources/test-resource-id/grants/test-policy-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Response.Data[0], http.StatusAccepted)
}

func TestPut_TypeFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/resources/test-resource-id/grants/test-policy-id", typeFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/resources/:resource_id/grants/:policy_id"),
		http.StatusBadRequest)
}

func TestPut_WithGrantFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/resources/test-resource-id/grants/test-policy-id", withGrantFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Attributes.WithGrant' Error:Field validation for 'WithGrant' failed on the 'required' tag", "/v1/resources/:resource_id/grants/:policy_id"),
		http.StatusBadRequest)
}

func TestPut_ServiceReturnsError_Returns500(t *testing.T) {
	mockService := mockService{
		Err: errors.New("some-test-error"),
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/resources/test-resource-id/grants/test-policy-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "some-test-error", "/v1/resources/:resource_id/grants/:policy_id"),
		http.StatusInternalServerError)
}
