package mocks

import (
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// MockResourceService struct to mock data returned from neo4j
type MockResourceService struct {
	mockGetResourceResponse     func() (resource.Response, error)
	mockGetResourceByIDResponse func(string) (resource.Element, error)
	mockCreateResourceResponse  func(*resource.Input) (resource.Element, error)
	mockDeleteResourceResponse  func(string) error
}
