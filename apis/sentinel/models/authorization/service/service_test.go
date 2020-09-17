package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/mocks"
	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/injection"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetAuthorizationForPrincipalResponse authorization.Output
	IsTargetOwnedByClientResponse        bool
	IsContextOwnedByClientResponse       bool
	IsPermissionOwnedByTenantResponse    bool
	Err                                  error
}

func (m mockRepository) GetAuthorizationForPrincipal(string, string, authorization.Input) (authorization.Output, error) {
	return m.GetAuthorizationForPrincipalResponse, m.Err
}

func (m mockRepository) IsTargetOwnedByClient(string, string, string) bool {
	return m.IsTargetOwnedByClientResponse
}

func (m mockRepository) IsContextOwnedByClient(string, string, string) bool {
	return m.IsContextOwnedByClientResponse
}

func (m mockRepository) IsPermissionOwnedByTenant(string, string, string) bool {
	return m.IsPermissionOwnedByTenantResponse
}

var errTest = errors.New("test-error")

func TestGetAuthorization_RepositoryReturnsData_ReturnsData(t *testing.T) {
	repository := mockRepository{GetAuthorizationForPrincipalResponse: testdata.Output}
	service := service.New(repository)

	output, err := service.GetAuthorizationForPrincipal("test-context-id", "test-principal-ID", authorization.Input{})
	assert.Equal(t, testdata.Output, output)
	assert.Nil(t, err)
}

func TestGetAuthorization_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	s := service.New(repository)

	_, err := s.GetAuthorizationForPrincipal("test-context-id", "test-principal-ID", authorization.Input{})
	assert.Equal(t, errTest, err)
}

func TestIsTargetOwnedByClient_RepositoryReturnsFalse_ReturnFalse(t *testing.T) {
	defer injection.Reset()
	injection.ExtractTenant = func(ctx context.Context) string { return "tenant" }
	injection.ExtractClaims = func(context.Context, string) string { return "clientID" }

	repository := mockRepository{IsTargetOwnedByClientResponse: false}
	s := service.New(repository)

	response := s.IsTargetOwnedByClient(mocks.Context{}, "target-id")
	assert.False(t, response)
}

func TestIsTargetOwnedByClient_RepositoryReturnsTrue_ReturnTrue(t *testing.T) {
	defer injection.Reset()
	injection.ExtractTenant = func(ctx context.Context) string { return "tenant" }
	injection.ExtractClaims = func(context.Context, string) string { return "clientID" }

	repository := mockRepository{IsTargetOwnedByClientResponse: true}
	s := service.New(repository)

	response := s.IsTargetOwnedByClient(mocks.Context{}, "target-id")
	assert.True(t, response)
}

func TestIsContextOwnedByTenant_RepositoryReturnsFalse_ReturnFalse(t *testing.T) {
	defer injection.Reset()
	injection.ExtractTenant = func(ctx context.Context) string { return "tenant" }
	injection.ExtractClaims = func(context.Context, string) string { return "clientID" }

	repository := mockRepository{IsContextOwnedByClientResponse: false}
	s := service.New(repository)

	response := s.IsContextOwnedByClient(mocks.Context{}, "test-context-id")
	assert.False(t, response)
}

func TestIsContextOwnedByTenant_RepositoryReturnsTrue_ReturnTrue(t *testing.T) {
	defer injection.Reset()
	injection.ExtractTenant = func(ctx context.Context) string { return "tenant" }
	injection.ExtractClaims = func(context.Context, string) string { return "clientID" }

	repository := mockRepository{IsContextOwnedByClientResponse: true}
	s := service.New(repository)

	response := s.IsContextOwnedByClient(mocks.Context{}, "test-context-id")
	assert.True(t, response)
}

func TestIsPermissionOwnedByTenant_RepositoryReturnsTrue_ReturnTrue(t *testing.T) {
	defer injection.Reset()
	injection.ExtractTenant = func(ctx context.Context) string { return "tenant" }
	injection.ExtractClaims = func(context.Context, string) string { return "clientID" }

	repository := mockRepository{IsPermissionOwnedByTenantResponse: true}
	s := service.New(repository)

	response := s.IsPermissionOwnedByTenant(mocks.Context{}, "test-permission-id")
	assert.True(t, response)
}

func TestIsPermissionOwnedByTenant_RepositoryReturnsFalse_ReturnFalse(t *testing.T) {
	defer injection.Reset()
	injection.ExtractTenant = func(ctx context.Context) string { return "tenant" }
	injection.ExtractClaims = func(context.Context, string) string { return "clientID" }

	repository := mockRepository{IsPermissionOwnedByTenantResponse: false}
	s := service.New(repository)

	response := s.IsPermissionOwnedByTenant(mocks.Context{}, "test-permission-id")
	assert.False(t, response)
}
