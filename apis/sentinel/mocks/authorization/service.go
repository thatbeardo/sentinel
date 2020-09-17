package mocks

import (
	"context"
	"testing"

	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/stretchr/testify/assert"
)

// AuthorizationService is a mock authorization service
type AuthorizationService struct {
	GetAuthorizationForPrincipalResponse authorization.Output
	IsTargetOwnedByClientResponse        bool
	IsPermissionOwnedByTenantResponse    bool
	IsContextOwnedByClientResponse       bool
	ExpectedInput                        authorization.Input
	Err                                  error
	T                                    *testing.T
}

// GetAuthorizationForPrincipal determines if the principal has permissions to the requestes target
func (m AuthorizationService) GetAuthorizationForPrincipal(principalID, contextID string, input authorization.Input) (output authorization.Output, err error) {
	assert.Equal(m.T, m.ExpectedInput, input)
	return m.GetAuthorizationForPrincipalResponse, m.Err
}

// IsTargetOwnedByClient determines if the target is owned by the tenant
func (m AuthorizationService) IsTargetOwnedByClient(context.Context, string) bool {
	return m.IsTargetOwnedByClientResponse
}

// IsContextOwnedByClient determines if the context under question is owned by the correct tenant
func (m AuthorizationService) IsContextOwnedByClient(context.Context, string) bool {
	return m.IsContextOwnedByClientResponse
}

// IsPermissionOwnedByTenant checks if the permission being updated is owned by the tenant
func (m AuthorizationService) IsPermissionOwnedByTenant(context.Context, string) bool {
	return m.IsPermissionOwnedByTenantResponse
}
