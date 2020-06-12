package service_test

import (
	"errors"
	"testing"

	authorization "github.com/bithippie/guard-my-app/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/service"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetAuthorizationForPrincipalResponse authorization.Output
	Err                                  error
}

func (m mockRepository) GetAuthorizationForPrincipal(string, authorization.Input) (authorization.Output, error) {
	return m.GetAuthorizationForPrincipalResponse, m.Err
}

var errTest = errors.New("test-error")

func TestGetAuthorization_RepositoryReturnsData_ReturnsData(t *testing.T) {
	repository := mockRepository{GetAuthorizationForPrincipalResponse: testdata.Output}
	service := service.New(repository)

	output, err := service.GetAuthorizationForPrincipal("test-principal-ID", authorization.Input{})
	assert.Equal(t, testdata.Output, output)
	assert.Nil(t, err)
}

func TestGetAuthorization_RepositoryReturnsError_ReturnsData(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.GetAuthorizationForPrincipal("test-principal-ID", authorization.Input{})
	assert.Equal(t, errTest, err)
}
