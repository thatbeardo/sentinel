package service_test

import (
	"errors"
	"testing"

	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetResponse     context.Output
	GetByIDResponse context.OutputDetails
	CreateResponse  context.OutputDetails
	UpdateResponse  context.OutputDetails
	DeleteResponse  context.Details

	GetByIDErr error
	Err        error
}

func (m mockRepository) Get() (context.Output, error) {
	return m.GetResponse, m.Err
}

func (m mockRepository) GetByID(id string) (context.OutputDetails, error) {
	return m.GetByIDResponse, m.Err
}

func (m mockRepository) Create(*context.Input) (context.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockRepository) Update(string, *context.Input) (context.OutputDetails, error) {
	if m.GetByIDErr != nil {
		return context.OutputDetails{}, m.GetByIDErr
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

	context, err := service.GetByID("test-id")
	assert.Equal(t, testdata.OutputDetails, context)
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

	context, err := service.Update("test-id", testdata.Input)
	assert.Equal(t, testdata.OutputDetails, context)
	assert.Nil(t, err)
}

func TestUpdate_contextNotFound_ReturnsNotFoundError(t *testing.T) {
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
