package service_test

import (
	"errors"
	"testing"

	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetResponse     policy.Output
	GetByIDResponse policy.OutputDetails
	CreateResponse  policy.OutputDetails
	UpdateResponse  policy.OutputDetails
	DeleteResponse  policy.Details

	GetByIDErr error
	Err        error
}

func (m mockRepository) Get() (policy.Output, error) {
	return m.GetResponse, m.Err
}

func (m mockRepository) GetByID(id string) (policy.OutputDetails, error) {
	return m.GetByIDResponse, m.Err
}

func (m mockRepository) Create(*policy.Input) (policy.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockRepository) Update(string, *policy.Input) (policy.OutputDetails, error) {
	if m.GetByIDErr != nil {
		return policy.OutputDetails{}, m.GetByIDErr
	}
	return m.UpdateResponse, m.Err
}

func (m mockRepository) Delete(string) error {
	return m.Err
}

var errTest = errors.New("test-error")

func TestGetByID_RepositoryReturnsError_ErrorReturned(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.GetByID("test-id")
	assert.Equal(t, errTest, err)
}

func TestGetByID_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{
		GetByIDResponse: testdata.OutputDetails,
	}
	service := service.New(repository)

	policy, err := service.GetByID("test-id")
	assert.Equal(t, testdata.OutputDetails, policy)
	assert.Nil(t, err)
}

func TestUpdate_RepositoryReturnsError_ErrorReturned(t *testing.T) {
	repository := mockRepository{
		Err: errTest,
	}
	service := service.New(repository)

	_, err := service.Update("test-id", testdata.Input)
	assert.Equal(t, errTest, err)
}

func TestUpdate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{UpdateResponse: testdata.OutputDetails}
	service := service.New(repository)

	policy, err := service.Update("test-id", testdata.Input)
	assert.Equal(t, testdata.OutputDetails, policy)
	assert.Nil(t, err)
}

func TestUpdate_PolicyNotFound_ReturnsNotFoundError(t *testing.T) {
	repository := mockRepository{GetByIDErr: errors.New("some-test-error")}
	service := service.New(repository)

	_, err := service.Update("test-id", testdata.Input)
	assert.Equal(t, errors.New("some-test-error"), err)
}

func TestDelete_RepositoryReturnsError_ErrorReturned(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	err := service.Delete("test-id")
	assert.Equal(t, errTest, err)
}

func TestDelete_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{DeleteResponse: testdata.Output.Data[0]}
	service := service.New(repository)

	err := service.Delete("test-id")
	assert.Nil(t, err)
}
