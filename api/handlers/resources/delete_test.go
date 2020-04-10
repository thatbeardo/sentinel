package resources_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/server"
	testutil "github.com/thatbeardo/go-sentinel/testutil/resources"
)

func TestDeleteResourcesOk(t *testing.T) {

	mockService := testutil.NewMockDeleteService(deleteResourceMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/resources/sample-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func deleteResourceMockResponseNoErrors(string) error {
	return nil
}
