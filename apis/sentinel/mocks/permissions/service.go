package mocks

import permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"

// PermissionService is the mock implementation of the permission service
type PermissionService struct {
	IsPermissionIdempotentResponse bool
}

// GetAllPermissionsForcontext is the mock implementation
func (p PermissionService) GetAllPermissionsForcontext(string) (permission.Output, error) {
	return permission.Output{}, nil
}

// GetAllPermissionsForcontextWithResource is a mock implementation
func (p PermissionService) GetAllPermissionsForcontextWithResource(string, string) (permission.Output, error) {
	return permission.Output{}, nil
}

// Create is a mock implementation to add a new permission
func (p PermissionService) Create(*permission.Input, string, string) (permission.OutputDetails, error) {
	return permission.OutputDetails{}, nil
}

// IsPermissionIdempotent is a mock implementation
func (p PermissionService) IsPermissionIdempotent(input *permission.Input, contextID, principalID string) (bool, error) {
	return p.IsPermissionIdempotentResponse, nil
}

// Update functions mock implementation
func (p PermissionService) Update(string, *permission.Input) (permission.OutputDetails, error) {
	return permission.OutputDetails{}, nil
}

// Delete functions mock implementation
func (p PermissionService) Delete(string) error {
	return nil
}
