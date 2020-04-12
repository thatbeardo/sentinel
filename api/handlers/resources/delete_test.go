package resources_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/mocks"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
)

func TestDeleteResourcesOk(t *testing.T) {

	mockService := &mocks.Service{}
	mockService.On("Delete", "test-id").Return(deleteResourceMockResponseNoErrors)
	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/resources/test-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func deleteResourceMockResponseNoErrors(string) error {
	return nil
}
