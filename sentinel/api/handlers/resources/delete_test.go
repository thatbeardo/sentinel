package resources_test

import (
	"net/http"
	"testing"

	handler "github.com/bithippie/guard-my-app/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/sentinel/api/handlers/resources"
	"github.com/bithippie/guard-my-app/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/service"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	resources.Routes(group, s)
	return r
}

func TestDeleteResourcesOk(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Delete", "test-id").Return(nil)
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/resources/test-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func TestDeleteResourcesServiceError(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Delete", "test-id").Return(models.ErrDatabase)
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/resources/test-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}
