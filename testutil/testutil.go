package testutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// MockResourceService struct to mock data returned from neo4j
type MockResourceService struct {
	mockGetResourceResponse    func() (resource.Response, error)
	mockCreateResourceResponse func(*resource.Input) (resource.Element, error)
}

// Get method returns mock test data
func (service *MockResourceService) Get() (resource.Response, error) {
	return service.mockGetResourceResponse()
}

// Create method creates a new node on the graph
func (service *MockResourceService) Create(resource *resource.Input) (resource.Element, error) {
	return service.mockCreateResourceResponse(resource)
}

// NewMockGetService mocks the service with desired data to be returned
func NewMockGetService(mock func() (resource.Response, error)) resource.Service {
	return &MockResourceService{
		mockGetResourceResponse: mock,
	}
}

// PerformRequest creates and returns an initialized ResponseRecorder
func PerformRequest(r http.Handler, method, path string) *http.Response {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

// ValidateResponse reads and asserts the response body content
func ValidateResponse(t *testing.T, response *http.Response, expected interface{}, code int) {
	assert.Equal(t, code, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	out, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(out), string(body))
}
