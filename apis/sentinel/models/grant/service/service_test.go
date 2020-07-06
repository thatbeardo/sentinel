package service_test

import (
	"errors"
	"testing"

	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	CreateResponse                            grant.OutputDetails
	GetPrincipalAndcontextForResourceResponse grant.Output
	Err                                       error
}

func (m mockRepository) Create(*grant.Input, string, string) (grant.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockRepository) GetPrincipalAndcontextForResource(string) (grant.Output, error) {
	return m.GetPrincipalAndcontextForResourceResponse, m.Err
}

var errTest = errors.New("test-error")

func TestCreate_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.Create(testdata.Input, "test-context-id", "test-target-id")
	assert.Equal(t, errTest, err)
}

func TestCreate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{CreateResponse: testdata.OutputDetails}
	service := service.New(repository)

	grant, err := service.Create(testdata.Input, "test-context-id", "test-target-id")
	assert.Equal(t, testdata.OutputDetails, grant)
	assert.Nil(t, err)
}

func TestGetAllPrincipalsAndContexts_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.GetPrincipalAndcontextForResource("test-context-id")
	assert.Equal(t, errTest, err)
}

func TestGetAllPrincipalsAndContexts_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{GetPrincipalAndcontextForResourceResponse: testdata.Output}
	service := service.New(repository)

	grant, err := service.GetPrincipalAndcontextForResource("test-context-id")
	assert.Equal(t, testdata.Output, grant)
	assert.Nil(t, err)
}
