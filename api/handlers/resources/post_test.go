package resources_test

import (
	"net/http"
	"testing"

	"github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
)

const payload = `{"data":{"type":"eer","attributes":{"source_id":"anothertest"}}}`

func TestPostResourcesOk(t *testing.T) {

	mockService := testutil.NewMockCreateService(createResourceMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/resources/", payload)
	defer cleanup()

	testutil.ValidateResponse(t, response, generateElement(), http.StatusAccepted)
}

func createResourceMockResponseNoErrors(*resource.Input) (resource.Element, error) {
	return generateElement(), nil
}
