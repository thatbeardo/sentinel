package policies_test

import (
	"net/http"
	"testing"

	handler "github.com/bithippie/guard-my-app/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/sentinel/api/handlers/policies"
	"github.com/bithippie/guard-my-app/sentinel/api/views"
	errors "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/service"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/testdata"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type mockService struct {
	CreateResponse outputs.Policy
	Err            error
}

func (m mockService) Create(*inputs.Payload) (outputs.Policy, error) {
	return m.CreateResponse, m.Err
}

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	policies.Routes(group, s)
	return r
}

const noErrors = `{"data":{"type":"policy","attributes":{"name":"valid-request"}}}`
const nameAbsentBadRequest = `{"data":{"type":"policy","attributes":{}}}`
const attributeAbsentBadRequest = `{"data":{"type":"policy"}}`

const typeAbsentBadRequest = `{"data":{"attributes":{"name":"valid-request"}}}`
const dataAbsentBadRequest = `{"dayta":{"type":"policy","attributes":{"name":"valid-request"}}}`

func TestPost_AllParametersPresent_Returns200(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Response.Data[0], http.StatusAccepted)
}

func TestPost_NameAttributeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", nameAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
}

func TestPost_AttributeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", nameAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
}

func TestPost_TypeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", typeAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
}

func TestPost_DataAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/policies/", dataAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data' Error:Field validation for 'Data' failed on the 'required' tag", "/v1/policies/"), http.StatusBadRequest)
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
