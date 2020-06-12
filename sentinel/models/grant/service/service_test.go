package service_test

import (
	"errors"
	"testing"

	grant "github.com/bithippie/guard-my-app/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/service"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	CreateResponse                           grant.OutputDetails
	GetPrincipalAndPolicyForResourceResponse grant.Output
	Err                                      error
}

func (m mockRepository) Create(*grant.Input, string, string) (grant.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockRepository) GetPrincipalAndPolicyForResource(string) (grant.Output, error) {
	return m.GetPrincipalAndPolicyForResourceResponse, m.Err
}

var errTest = errors.New("test-error")

func TestCreate_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.NewService(repository)

	_, err := service.Create(testdata.Input, "test-policy-id", "test-target-id")
	assert.Equal(t, errTest, err)
}

func TestCreate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{CreateResponse: testdata.OutputDetails}
	service := service.NewService(repository)

	grant, err := service.Create(testdata.Input, "test-policy-id", "test-target-id")
	assert.Equal(t, testdata.OutputDetails, grant)
	assert.Nil(t, err)
}

func TestGetAllPrincipalsAndPolicies_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.NewService(repository)

	_, err := service.GetPrincipalAndPolicyForResource("test-policy-id")
	assert.Equal(t, errTest, err)
}

func TestGetAllPrincipalsAndPolicies_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{GetPrincipalAndPolicyForResourceResponse: testdata.Output}
	service := service.NewService(repository)

	grant, err := service.GetPrincipalAndPolicyForResource("test-policy-id")
	assert.Equal(t, testdata.Output, grant)
	assert.Nil(t, err)
}
