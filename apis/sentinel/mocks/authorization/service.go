package mocks

import (
	"testing"

	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/stretchr/testify/assert"
)

// AuthorizationService is a mock authorization service
type AuthorizationService struct {
	GetAuthorizationForPrincipalResponse authorization.Output
	IsTargetOwnedByTenantResponse        bool
	IsPermissionOwnedByTenantResponse    bool
	IsPolicyOwnedByTenantResponse        bool
	ExpectedInput                        authorization.Input
	Err                                  error
	T                                    *testing.T
}

// GetAuthorizationForPrincipal determines if the principal has permissions to the requestes target
func (m AuthorizationService) GetAuthorizationForPrincipal(principalID string, input authorization.Input) (output authorization.Output, err error) {
	assert.Equal(m.T, m.ExpectedInput, input)
	return m.GetAuthorizationForPrincipalResponse, m.Err
}

// IsTargetOwnedByTenant determines if the target is owned by the tenant
func (m AuthorizationService) IsTargetOwnedByTenant(string, string) bool {
	return m.IsTargetOwnedByTenantResponse
}

// IsPolicyOwnedByTenant determines if the policy under question is owned by the correct tenant
func (m AuthorizationService) IsPolicyOwnedByTenant(string, string) bool {
	return m.IsPolicyOwnedByTenantResponse
}

// IsPermissionOwnedByTenant checks if the permission being updated is owned by the tenant
func (m AuthorizationService) IsPermissionOwnedByTenant(string, string) bool {
	return m.IsPermissionOwnedByTenantResponse
}
