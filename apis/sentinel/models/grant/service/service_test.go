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
	ExpectedContextID                         string
	ExpectedPrincipalID                       string
	GrantExistsResponse                       bool
	Err                                       error
	T                                         *testing.T
}

func (m mockRepository) Create(*grant.Input, string, string) (grant.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockRepository) GetPrincipalAndcontextForResource(string) (grant.Output, error) {
	return m.GetPrincipalAndcontextForResourceResponse, m.Err
}

func (m mockRepository) GrantExists(contextID, principalID string) (bool, error) {
	assert.Equal(m.T, m.ExpectedContextID, contextID)
	assert.Equal(m.T, m.ExpectedPrincipalID, principalID)
	return m.GrantExistsResponse, m.Err
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

func TestGrantExists_RepositoryReturnsTrue_ReturnTrue(t *testing.T) {
	repository := mockRepository{
		GrantExistsResponse: true,
		ExpectedContextID:   "test-context-id",
		ExpectedPrincipalID: "test-principal-id",
		T:                   t,
	}
	service := service.New(repository)

	grantExists, err := service.GrantExists("test-context-id", "test-principal-id")
	assert.Equal(t, grantExists, true)
	assert.Nil(t, err, nil)
}

func TestGrantExists_RepositoryReturnsError_ReturnError(t *testing.T) {
	repository := mockRepository{
		Err:                 errTest,
		ExpectedContextID:   "test-context-id",
		ExpectedPrincipalID: "test-principal-id",
		T:                   t,
	}
	service := service.New(repository)

	grantExists, err := service.GrantExists("test-context-id", "test-principal-id")
	assert.Equal(t, grantExists, false)
	assert.Equal(t, err, errTest)
}
