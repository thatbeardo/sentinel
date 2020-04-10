package resources_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/thatbeardo/go-sentinel/api/views"
	"github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/server"
	testutil "github.com/thatbeardo/go-sentinel/testutil/resources"
)

func TestInvalidPath(t *testing.T) {

	mockService := testutil.NewMockGetService(getResourceMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resurces/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateError("", "query-parameter-todo", "Path not found", http.StatusNotFound), http.StatusNotFound)
}

func TestGetResourcesOk(t *testing.T) {

	mockService := testutil.NewMockGetService(getResourceMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateResponse(), http.StatusOK)
}

func TestGetResourcesDatabaseError(t *testing.T) {

	mockService := testutil.NewMockGetService(getReourceMockResponse500)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateError("/v1/resources/", "query-parameter-todo", "Database Error", http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func TestGetResourceByIDOk(t *testing.T) {

	mockService := testutil.NewMockGetByIDService(getResourceByIdMockResponseNoErrors)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/sample-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateElement(), http.StatusOK)
}

func TestGetResourceByIDNoResourceFound(t *testing.T) {

	mockService := testutil.NewMockGetByIDService(getResourceByIdMockResponseNoResource)

	router := server.SetupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/resources/sample-id", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, generateError("/v1/resources/:id", "query-parameter-todo", "Document not found", http.StatusNotFound), http.StatusNotFound)
}

func getResourceMockResponseNoErrors() (resource.Response, error) {
	return generateResponse(), nil
}

func getResourceByIdMockResponseNoErrors(string) (resource.Element, error) {
	return generateElement(), nil
}

func getResourceByIdMockResponseNoResource(string) (resource.Element, error) {
	return generateElement(), errors.New("Document not found")
}

func getReourceMockResponse500() (resource.Response, error) {
	return resource.Response{}, errors.New("Database Error")
}

func generateResponse() resource.Response {
	element := generateElement()
	return resource.Response{Data: []resource.Element{element}}
}

func generateElement() resource.Element {
	policies := resource.Policies{Data: []resource.Identifier{}}
	parent := resource.Parent{Data: resource.Identifier{Type: "parent-resource", ID: "parent-resource-id"}}
	relationships := resource.Relationships{Parent: parent, Policies: policies}
	userResource := resource.Resource{Name: "my-resource", SourceID: "my-source-id"}
	return resource.Element{Relationships: relationships, Attributes: userResource, Type: "resource", ID: "uuid"}
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
