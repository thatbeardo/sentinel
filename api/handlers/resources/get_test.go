package resources_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/thatbeardo/go-sentinel/api/views"
	"github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
)

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

	testutil.ValidateResponse(t, response, generateError("/v1/resources/", "query-parameter-todo", "Database Error"),
		http.StatusInternalServerError)
}

func getResourceMockResponseNoErrors() (resource.Response, error) {
	return generateResponse(), nil
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

func generateError(pointer string, parameter string, detail string) views.ErrView {
	source := views.Source{
		Pointer:   pointer,
		Parameter: parameter,
	}
	return views.ErrView{
		ID:     "error-id-todo",
		Status: 500,
		Source: source,
		Detail: detail,
	}
}
