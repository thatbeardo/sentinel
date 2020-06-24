package service_test

import (
	"errors"
	"testing"

	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetAuthorizationForPrincipalResponse authorization.Output
	IsTargetOwnedByTenantResponse        bool
	Err                                  error
}

func (m mockRepository) GetAuthorizationForPrincipal(string, authorization.Input) (authorization.Output, error) {
	return m.GetAuthorizationForPrincipalResponse, m.Err
}

func (m mockRepository) IsTargetOwnedByTenant(string, string) bool {
	return m.IsTargetOwnedByTenantResponse
}

var errTest = errors.New("test-error")

func TestGetAuthorization_RepositoryReturnsData_ReturnsData(t *testing.T) {
	repository := mockRepository{GetAuthorizationForPrincipalResponse: testdata.Output}
	service := service.New(repository)

	output, err := service.GetAuthorizationForPrincipal("test-principal-ID", authorization.Input{})
	assert.Equal(t, testdata.Output, output)
	assert.Nil(t, err)
}

func TestGetAuthorization_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service.SetRepository(repository)
	s := service.New(repository)

	_, err := s.GetAuthorizationForPrincipal("test-principal-ID", authorization.Input{})
	assert.Equal(t, errTest, err)
}

func TestIsTargetOwnedByTenant_RepositoryReturnsFalse_ReturnFalse(t *testing.T) {
	repository := mockRepository{IsTargetOwnedByTenantResponse: false}
	service.SetRepository(repository)
	s := service.New(repository)

	response := s.IsTargetOwnedByTenant("target-id", "tenant-id")
	assert.False(t, response)
}

func TestIsTargetOwnedByTenant_RepositoryReturnsTrue_ReturnTrue(t *testing.T) {
	repository := mockRepository{IsTargetOwnedByTenantResponse: true}
	service.SetRepository(repository)
	s := service.NewWithoutRepository()

	response := s.IsTargetOwnedByTenant("target-id", "tenant-id")
	assert.True(t, response)
}
