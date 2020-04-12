package resources_test

import (
	"net/http"
	"testing"

	"github.com/thatbeardo/go-sentinel/api/views"
	"github.com/thatbeardo/go-sentinel/mocks"

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

func TestPostResourcesOk(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Create", m.AnythingOfType("*resource.Input")).Return(createResourceMockResponseNoErrors())

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

func createResourceMockResponseNoErrors() (resource.Element, error) {
	return generateElement(), nil
}
