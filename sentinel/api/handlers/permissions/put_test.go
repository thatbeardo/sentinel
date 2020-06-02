package permissions_test

import (
	"net/http"
	"testing"

	handler "github.com/bithippie/guard-my-app/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/sentinel/api/handlers/permissions"
	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/service"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/testdata"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type mockService struct {
	CreateResponse outputs.Permission
	Err            error
}

func (m mockService) Create(*inputs.Payload, string, string) (outputs.Permission, error) {
	return m.CreateResponse, m.Err
}

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	permissions.Routes(group, s)
	return r
}

const noErrors = `{"data":{"type":"permission","attributes":{"name":"resource:read","permitted":"allow"}}}`
const typeFieldAbsent = `{"data":{"attributes":{"name":"resource:read","permitted":"allow"}}}`
const nameFieldAbsent = `{"data":{"type":"permission","attributes":{"permitted":"allow"}}}`
const permittedFieldAbsent = `{"data":{"type":"permission","attributes":{"name":"resource:read"}}}`

func TestPut_AllParametersPresent_Returns200(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/policies/test-policy-id/resources/test-target-id/permissions", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.Response.Data[0], http.StatusAccepted)
}

func TestPut_TypeFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/policies/test-policy-id/resources/test-target-id/permissions", typeFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/policies/:policy_id/resources/:resource_id/permissions"),
		http.StatusBadRequest)
}

func TestPut_NameFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/policies/test-policy-id/resources/test-target-id/permissions", nameFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag", "/v1/policies/:policy_id/resources/:resource_id/permissions"),
		http.StatusBadRequest)
}

func TestPut_PermittedFieldAbsent_ReturnsBadRequest(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.Response.Data[0],
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PUT", "/v1/policies/test-policy-id/resources/test-target-id/permissions", permittedFieldAbsent)
	defer cleanup()

	testutil.ValidateResponse(t,
		response,
		views.GenerateErrorResponse(response.StatusCode, "Key: 'Payload.Data.Attributes.Permitted' Error:Field validation for 'Permitted' failed on the 'required' tag", "/v1/policies/:policy_id/resources/:resource_id/permissions"),
		http.StatusBadRequest)
}
