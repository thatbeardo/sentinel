package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/thatbeardo/go-sentinel/models/resource"
)

// MockResourceService struct to mock data returned from neo4j
type MockResourceService struct {
	mockGetResource    func() (resource.Response, error)
	mockCreateResource func(*resource.Input) (resource.Element, error)
}

// Get method returns mock test data
func (service *MockResourceService) Get() (resource.Response, error) {
	return service.mockGetResource()
}

// Create method creates a new node on the graph
func (service *MockResourceService) Create(resource *resource.Input) (resource.Element, error) {
	return service.mockCreateResource(resource)
}

// PerformRequest creates and returns an initialized ResponseRecorder
func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// GetMockService mocks resource service
func GetMockService() resource.Service {
	return &MockResourceService{
		mockGetResource:    getResourceMockData,
		mockCreateResource: createResourceMockData,
	}
}

func getResourceMockData() (resource.Response, error) {
	return resource.Response{Data: []resource.Element{}}, nil
}

func createResourceMockData(*resource.Input) (resource.Element, error) {
	return resource.Element{}, nil
}
