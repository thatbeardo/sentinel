package testutil

import "github.com/thatbeardo/go-sentinel/models/resource"

// Delete function to delete a node from the graph
func (service *MockResourceService) Delete(id string) error {
	return service.mockDeleteResourceResponse(id)
}

// NewMockDeleteService mocks the service with desired data to be returned when POST is called
func NewMockDeleteService(mock func(string) error) resource.Service {
	return &MockResourceService{
		mockDeleteResourceResponse: mock,
	}
}
