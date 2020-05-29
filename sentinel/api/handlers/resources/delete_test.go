package resources_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/bithippie/go-sentinel/sentinel/mocks"
	models "github.com/bithippie/go-sentinel/sentinel/models"
	"github.com/bithippie/go-sentinel/sentinel/server"
	"github.com/bithippie/go-sentinel/sentinel/testutil"
)

func TestDeleteResourcesOk(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Delete", "test-id").Return(nil)
	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/resources/test-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func TestDeleteResourcesServiceError(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Delete", "test-id").Return(models.ErrDatabase)
	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/resources/test-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}
