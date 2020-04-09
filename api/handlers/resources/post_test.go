package resources_test

import (
	"net/http"
	"testing"

	"github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
)

const noErrors = `{"data":{"type":"resource","attributes":{"source_id":"much-required"}}}`
const sourceIdAbsent = `{"data":{"type":"resource","attributes":{"someField":"not-much-required"}}}`

func TestPostResourcesOk(t *testing.T) {

	mockService := testutil.NewMockCreateService(createResourceMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, generateElement(), http.StatusAccepted)
}

func TestPostResourcesSourceIdAbsent(t *testing.T) {

	mockService := testutil.NewMockCreateService(createResourceMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", sourceIdAbsent)
	defer cleanup()

	testutil.ValidateResponse(t, response, "Key: 'Input.Data.Attributes.SourceID' Error:Field validation for 'SourceID' failed on the 'required' tag", http.StatusBadRequest)
}

func createResourceMockResponseNoErrors(*resource.Input) (resource.Element, error) {
	return generateElement(), nil
}
