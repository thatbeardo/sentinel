package resources_test

import (
	"net/http"
	"testing"

	m "github.com/stretchr/testify/mock"
	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/mocks"
	"github.com/bithippie/guard-my-app/sentinel/server"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
)

func TestPatchResourceOk(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Update", "test-id", m.AnythingOfType("*entity.Input")).Return(createResourceNoErrors())

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/resources/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, generateElement(), http.StatusAccepted)
}

func TestPatchResourceDatabaseError(t *testing.T) {
	mockService := &mocks.Service{}
	mockService.On("Update", "test-id", m.AnythingOfType("*entity.Input")).Return(databaseError())

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/resources/test-id", noErrors)
	defer cleanup()

	testutil.ValidateResponse(t, response, generateError("/v1/resources/:id", "query-parameter-todo", "Database Error", http.StatusInternalServerError), http.StatusInternalServerError)
}

func TestPatchResourcesSourceIdBlank(t *testing.T) {
	mockService := &mocks.Service{}

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "PATCH", "/v1/resources/test-id", sourceIdBlank)
	defer cleanup()

	testutil.ValidateResponse(t, response, views.GenerateErrorResponse(http.StatusBadRequest, "Key: 'Input.Data.Attributes.SourceID' Error:Field validation for 'SourceID' failed on the 'required' tag", "/v1/resources/:id"), http.StatusBadRequest)
}
