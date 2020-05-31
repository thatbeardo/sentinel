package service_test

import (
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/service"
	testdata "github.com/bithippie/guard-my-app/sentinel/models/policy/test-data"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	CreateResponse outputs.Policy
	Err            error
}

func (m mockRepository) Create(*inputs.Payload) (outputs.Policy, error) {
	return m.CreateResponse, m.Err
}

var errTest = errors.New("test-error")

func TestCreate_RepositoryReturnsError_ErrorReturned(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.NewService(repository)

	_, err := service.Create(testdata.Payload)
	assert.Equal(t, errTest, err)
}

func TestCreate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{CreateResponse: testdata.Response.Data[0]}
	service := service.NewService(repository)

	policy, err := service.Create(testdata.Payload)
	assert.Equal(t, testdata.Response.Data[0], policy)
	assert.Nil(t, err)
}
