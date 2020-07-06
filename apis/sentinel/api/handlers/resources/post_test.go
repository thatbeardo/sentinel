package resources_test

import (
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	errors "github.com/bithippie/guard-my-app/apis/sentinel/models"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	contextTestData "github.com/bithippie/guard-my-app/apis/sentinel/models/context/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

const noErrors = `{"data":{"type":"resource","attributes":{"source_id":"much-required"}}}`
const sourceIdAbsent = `{"data":{"type":"resource","attributes":{"someField":"not-much-required"}}}`
const sourceIdBlank = `{"data":{"type":"resource","attributes":{"source_id":""}}}`
const typeAbsent = `{"data":{"typo":"resource","attributes":{"source_id":"not-much-required"}}}`
const typeBlank = `{"data":{"type":"","attributes":{"source_id":"much-required"}}}`

const relationshipsEmptyPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{}}}`
const relationshipsParentDataAbsentPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{"parent":{"type":"resource"}}}}`
const parentDataIDAbsentPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{"parent":{"data":{"type":"resource"}}}}}`
const parentDataTypeAbsentPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{"parent":{"data":{"id":"test-id"}}}}}`

const associatecontextNoErrors = `{"data":{"type":"context","attributes":{"name":"valid-request"}}}`
const nameAbsentBadRequest = `{"data":{"type":"context","attributes":{}}}`
const attributeAbsentBadRequest = `{"data":{"type":"context"}}`

const typeAbsentBadRequest = `{"data":{"attributes":{"name":"valid-request"}}}`
const dataAbsentBadRequest = `{"dayta":{"type":"context","attributes":{"name":"valid-request"}}}`

func TestPostResourcesOk(t *testing.T) {
	mockService := mockService{
		CreateResponse: testdata.ModificationResult,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, testdata.ModificationResult, http.StatusAccepted)
}

func TestPostResourcesSourceIdBlank(t *testing.T) {

	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", sourceIdBlank)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.SourceID' Error:Field validation for 'SourceID' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourcesSourceIdAbsent(t *testing.T) {

	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", sourceIdAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.SourceID' Error:Field validation for 'SourceID' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourcesTypeBlank(t *testing.T) {

	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", typeBlank)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourcesTypeAbsent(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", typeAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceEmptyRelationships(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", relationshipsEmptyPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent' Error:Field validation for 'Parent' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentDataAbsent(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", relationshipsParentDataAbsentPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent.Data' Error:Field validation for 'Data' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentIdAbsent(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", parentDataIDAbsentPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent.Data.ID' Error:Field validation for 'ID' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentTypeAbsent(t *testing.T) {
	mockService := mockService{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", parentDataTypeAbsentPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentAbsentInDatabase(t *testing.T) {
	mockService := mockService{
		CreateErr: models.ErrNotFound,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusNotFound, "Data not found", "/v1/resources/"), http.StatusNotFound)
}

func TestPost_AllParametersPresent_Returns200(t *testing.T) {
	mockService := mockService{
		AssociateResponse: contextTestData.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/test-id/contexts", associatecontextNoErrors)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		contextTestData.OutputDetails,
		http.StatusAccepted)
}

func TestPost_NameAttributeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		AssociateResponse: contextTestData.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/test-id/contexts", nameAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			response.StatusCode,
			"Key: 'Input.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag",
			"/v1/resources/:id/contexts"),
		http.StatusBadRequest)
}

func TestPost_AttributeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		AssociateResponse: contextTestData.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/test-id/contexts", nameAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			response.StatusCode,
			"Key: 'Input.Data.Attributes.Name' Error:Field validation for 'Name' failed on the 'required' tag",
			"/v1/resources/:id/contexts"),
		http.StatusBadRequest)
}

func TestPost_TypeAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		AssociateResponse: contextTestData.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/test-id/contexts", typeAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			response.StatusCode,
			"Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag",
			"/v1/resources/:id/contexts"),
		http.StatusBadRequest)
}

func TestPost_DataAbsent_Returns400(t *testing.T) {
	mockService := mockService{
		AssociateResponse: contextTestData.OutputDetails,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/test-id/contexts", dataAbsentBadRequest)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			response.StatusCode,
			"Key: 'Input.Data' Error:Field validation for 'Data' failed on the 'required' tag",
			"/v1/resources/:id/contexts"),
		http.StatusBadRequest)
}

func TestPost_ServiceReturnsError_Returns500(t *testing.T) {
	mockService := mockService{
		AssociateErr: errors.ErrDatabase,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/test-id/contexts", associatecontextNoErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(response.StatusCode, "Database Error", "/v1/resources/:id/contexts"), http.StatusInternalServerError)
}
