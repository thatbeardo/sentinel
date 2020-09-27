package mocks

import (
	"testing"

	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
)

// GrantService is a mocked struct used for middleware testing
type GrantService struct {
	ExpectedContextID   string
	ExpectedPrincipalID string
	GrantExistsResponse bool
	Err                 error
	T                   *testing.T
}

// Create is a mocked function to add a grant
func (g GrantService) Create(*grant.Input, string, string) (grant.OutputDetails, error) {
	return grant.OutputDetails{}, nil
}

// GetPrincipalAndcontextForResource is a mocked implementation of the service
func (g GrantService) GetPrincipalAndcontextForResource(string) (grant.Output, error) {
	return grant.Output{}, nil
}

// GrantExists is used to determine duplicate grants being created
func (g GrantService) GrantExists(contextID, principalID string) (bool, error) {
	return g.GrantExistsResponse, g.Err
}
