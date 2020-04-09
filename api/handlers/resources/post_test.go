package resources_test

import (
	"net/http"
	"testing"

	"github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
)

func TestPostResourcesOk(t *testing.T) {

	mockService := testutil.NewMockCreateService(createResourceMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response := testutil.PerformRequest(router, "POST", "/v1/resources/", createInput())

	testutil.ValidateResponse(t, response, generateElement(), http.StatusAccepted)
}

func createResourceMockResponseNoErrors(*resource.Input) (resource.Element, error) {
	return generateElement(), nil
}

func createInput() string {
	return `{"input":{"data":{"attributes":{"source_id":"asd"}}},"data":{"attributes":{"name":"string","source_id":"string"},"relationships":{"parent":{"data":{"id":"string","type":"policy"}}},"type":"string"}}`
}
