package testutil

import "github.com/thatbeardo/go-sentinel/models/resource"

// Create method creates a new node on the graph
func (service *MockResourceService) Create(resource *resource.Input) (resource.Element, error) {
	return service.mockCreateResourceResponse(resource)
}

// NewMockCreateService mocks the service with desired data to be returned when POST is called
func NewMockCreateService(mock func(*resource.Input) (resource.Element, error)) resource.Service {
	return &MockResourceService{
		mockCreateResourceResponse: mock,
	}
}
