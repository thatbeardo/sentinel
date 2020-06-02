package service_test

import (
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/sentinel/models/grant/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/service"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	CreateResponse outputs.Grant
	Err            error
}

func (m mockRepository) Create(*inputs.Payload, string, string) (outputs.Grant, error) {
	return m.CreateResponse, m.Err
}

var errTest = errors.New("test-error")

func TestCreate_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.NewService(repository)

	_, err := service.Create(testdata.Payload, "test-policy-id", "test-target-id")
	assert.Equal(t, errTest, err)
}

func TestCreate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{CreateResponse: testdata.Response.Data[0]}
	service := service.NewService(repository)

	permission, err := service.Create(testdata.Payload, "test-policy-id", "test-target-id")
	assert.Equal(t, testdata.Response.Data[0], permission)
	assert.Nil(t, err)
}
