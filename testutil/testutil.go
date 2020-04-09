package testutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

// NewMockCreateService mocks the service with desired data to be returned when POST is called
func NewMockCreateService(mock func(*resource.Input) (resource.Element, error)) resource.Service {
	return &MockResourceService{
		mockCreateResourceResponse: mock,
	}
}

// PerformRequest creates and returns an initialized ResponseRecorder
func PerformRequest(r http.Handler, method, path string, body string) (*http.Response, func() error) {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)
	response := recorder.Result()
	return response, response.Body.Close
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
