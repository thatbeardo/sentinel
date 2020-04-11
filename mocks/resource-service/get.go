package mocks

import (
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// NewMockGetService mocks the service with desired data to be returned
func NewMockGetService(mock func() (resource.Response, error)) resource.Service {
	return &MockResourceService{
		mockGetResourceResponse: mock,
	}
}

// NewMockGetByIDService mocks the service with desired data to be returned
func NewMockGetByIDService(mock func(string) (resource.Element, error)) resource.Service {
	return &MockResourceService{
		mockGetResourceByIDResponse: mock,
	}
}

// Get method returns mock test data
func (service *MockResourceService) Get() (resource.Response, error) {
	return service.mockGetResourceResponse()
}

// GetByID method gets a resource based on its ID
func (service *MockResourceService) GetByID(id string) (resource.Element, error) {
	return service.mockGetResourceByIDResponse(id)
}
