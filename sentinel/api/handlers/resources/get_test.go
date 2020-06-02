package resources_test

import (
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/sentinel/api/views"
	"github.com/bithippie/guard-my-app/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/sentinel/models"
	entity "github.com/bithippie/guard-my-app/sentinel/models/resource"
	"github.com/bithippie/guard-my-app/sentinel/testutil"
)

func TestInvalidPath(t *testing.T) {

	mockService := &mocks.Service{}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/invalid-path/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateError("", "query-parameter-todo", "Path not found", http.StatusNotFound), http.StatusNotFound)
}

func TestGetResourcesOk(t *testing.T) {

	mockService := &mocks.Service{}
	mockService.On("Get").Return(getResourceMockResponseNoErrors())

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateResponse(), http.StatusOK)
}

func TestGetResourcesDatabaseError(t *testing.T) {

	mockService := &mocks.Service{}
	mockService.On("Get").Return(getResourceReturns500())

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateError("/v1/resources/", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func TestGetResourceByIDOk(t *testing.T) {

	mockService := &mocks.Service{}
	mockService.On("GetByID", "test-id").Return(getResourceByIdMockResponseNoErrors())

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateElement(), http.StatusOK)
}

func TestGetResourceByIDNoResourceFound(t *testing.T) {

	mockService := &mocks.Service{}
	mockService.On("GetByID", "test-id").Return(getResourceByIdNoResource())

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/test-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateError("/v1/resources/:id", "query-parameter-todo", "Data not found", http.StatusNotFound), http.StatusNotFound)
}

func getResourceMockResponseNoErrors() (entity.Response, error) {
	return generateResponse(), nil
}

func getResourceByIdMockResponseNoErrors() (entity.Element, error) {
	return generateElement(), nil
}

func getResourceByIdNoResource() (entity.Element, error) {
	return generateElement(), models.ErrNotFound
}

func getResourceReturns500() (entity.Response, error) {
	return entity.Response{}, models.ErrDatabase
}

func generateResponse() entity.Response {
	element := generateElement()
	return entity.Response{Data: []entity.Element{element}}
}

func generateElement() entity.Element {
	policies := &entity.Policies{Data: []*entity.Identifier{}}
	parent := &entity.Parent{Data: &entity.Identifier{Type: "parent-resource", ID: "parent-resource-id"}}
	relationships := &entity.Relationships{Parent: parent, Policies: policies}
	userResource := &entity.Resource{Name: "my-resource", SourceID: "my-source-id"}
	return entity.Element{Relationships: relationships, Attributes: userResource, Type: "resource", ID: "uuid"}
}

func generateError(pointer string, parameter string, detail string, code int) views.ErrView {
	source := views.Source{
		Pointer:   pointer,
		Parameter: parameter,
	}
	return views.ErrView{
		ID:     "error-id-todo",
		Status: code,
		Source: source,
		Detail: detail,
	}
}
