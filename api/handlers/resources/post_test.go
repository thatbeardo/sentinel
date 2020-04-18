package resources_test

import (
	"net/http"
	"testing"

	"github.com/thatbeardo/go-sentinel/api/views"
	"github.com/thatbeardo/go-sentinel/mocks"
	models "github.com/thatbeardo/go-sentinel/models"

	m "github.com/stretchr/testify/mock"
	"github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
)

const noErrors = `{"data":{"type":"resource","attributes":{"source_id":"much-required"}}}`
const sourceIdAbsent = `{"data":{"type":"resource","attributes":{"someField":"not-much-required"}}}`
const sourceIdBlank = `{"data":{"type":"resource","attributes":{"source_id":""}}}`
const typeAbsent = `{"data":{"typoo":"resource","attributes":{"source_id":"not-much-required"}}}`
const typeBlank = `{"data":{"type":"","attributes":{"source_id":"much-required"}}}`

const relationshipsEmptyPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{}}}`
const relationshipsParentDataAbsentPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{"parent":{"type":"resource"}}}}`
const parentDataIDAbsentPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{"parent":{"data":{"type":"resource"}}}}}`
const parentDataTypeAbsentPayload = `{"data":{"type":"resource","attributes":{"source_id":"test-id"},"relationships":{"parent":{"data":{"id":"test-id"}}}}}`

func TestPostResourcesOk(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Create", m.AnythingOfType("*resource.Input")).Return(createResourceNoErrors())

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, generateElement(), http.StatusAccepted)
}

func TestPostResourcesSourceIdBlank(t *testing.T) {

	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", sourceIdBlank)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.SourceID' Error:Field validation for 'SourceID' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourcesSourceIdAbsent(t *testing.T) {

	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", sourceIdAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.SourceID' Error:Field validation for 'SourceID' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourcesTypeBlank(t *testing.T) {

	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", typeBlank)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourcesTypeAbsent(t *testing.T) {
	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", typeAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceEmptyRelationships(t *testing.T) {
	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", relationshipsEmptyPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent' Error:Field validation for 'Parent' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentDataAbsent(t *testing.T) {
	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", relationshipsParentDataAbsentPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent.Data' Error:Field validation for 'Data' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentIdAbsent(t *testing.T) {
	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", parentDataIDAbsentPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent.Data.ID' Error:Field validation for 'ID' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentTypeAbsent(t *testing.T) {
	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", parentDataTypeAbsentPayload)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Relationships.Parent.Data.Type' Error:Field validation for 'Type' failed on the 'required' tag", "/v1/resources/"), http.StatusBadRequest)
}

func TestPostResourceParentAbsentInDatabase(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Create", m.AnythingOfType("*resource.Input")).Return(createResourceParentNotFound())

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusNotFound, "Data not found", "/v1/resources/"), http.StatusNotFound)
}

func createResourceNoErrors() (resource.Element, error) {
	return generateElement(), nil
}

func createResourceParentNotFound() (resource.Element, error) {
	return resource.Element{}, models.ErrNotFound
}

func databaseError() (resource.Element, error) {
	return resource.Element{}, models.ErrDatabase
}
