package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/thatbeardo/go-sentinel/models/resource"
)

// MockResourceService struct to mock data returned from neo4j
type MockResourceService struct {
	mockGetResource func() (*resource.Resource, error)
}

// Get method returns mock test data
func (service *MockResourceService) Get() (*resource.Resource, error) {
	return service.mockGetResource()
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
	return &MockResourceService{mockGetResource: getResourceMockData}
}

func getResourceMockData() (*resource.Resource, error) {
	return &resource.Resource{Name: "Harshil", SourceID: "Mavain"}, nil
}
