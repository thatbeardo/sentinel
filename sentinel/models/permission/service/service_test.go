package service_test

import (
	"errors"
	"testing"

	permission "github.com/bithippie/guard-my-app/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/service"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetAllPermissionsForPolicyResponse             permission.Output
	GetAllPermissionsForPolicyResponseWithResource permission.Output
	CreateResponse                                 permission.OutputDetails
	UpdateResponse                                 permission.OutputDetails
	Err                                            error
}

func (m mockRepository) Create(*permission.Input, string, string) (permission.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockRepository) GetAllPermissionsForPolicy(policyID string) (permission.Output, error) {
	return m.GetAllPermissionsForPolicyResponse, m.Err
}

func (m mockRepository) GetAllPermissionsForPolicyWithResource(policyID string, resourceID string) (permission.Output, error) {
	return m.GetAllPermissionsForPolicyResponseWithResource, m.Err
}

func (m mockRepository) Update(string, *permission.Input) (permission.OutputDetails, error) {
	return m.UpdateResponse, m.Err
}

func (m mockRepository) Delete(string) error {
	return m.Err
}

var errTest = errors.New("test-error")

func TestCreate_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.Create(testdata.Input, "test-policy-id", "test-target-id")
	assert.Equal(t, errTest, err)
}

func TestCreate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{CreateResponse: testdata.OutputDetails}
	service := service.New(repository)

	permission, err := service.Create(testdata.Input, "test-policy-id", "test-target-id")
	assert.Equal(t, testdata.OutputDetails, permission)
	assert.Nil(t, err)
}

func TestGetAllPermissions_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.GetAllPermissionsForPolicy("test-policy-id")
	assert.Equal(t, errTest, err)
}

func TestGetAllPermissions_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{GetAllPermissionsForPolicyResponse: testdata.Output}
	service := service.New(repository)

	permission, err := service.GetAllPermissionsForPolicy("test-policy-id")
	assert.Equal(t, testdata.Output, permission)
	assert.Nil(t, err)
}

func TestGetAllPermissionsOfResource_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.GetAllPermissionsForPolicyWithResource("test-policy-id", "test-resource-id")
	assert.Equal(t, errTest, err)
}

func TestGetAllPermissionsOfResource_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{GetAllPermissionsForPolicyResponseWithResource: testdata.Output}
	service := service.New(repository)

	permission, err := service.GetAllPermissionsForPolicyWithResource("test-policy-id", "test-resource-id")
	assert.Equal(t, testdata.Output, permission)
	assert.Nil(t, err)
}

func TestUpdate_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.Update("test-policy-id", testdata.Input)
	assert.Equal(t, errTest, err)
}

func TestUpdate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{UpdateResponse: testdata.OutputDetails}
	service := service.New(repository)

	permission, err := service.Update("test-policy-id", testdata.Input)
	assert.Equal(t, testdata.OutputDetails, permission)
	assert.Nil(t, err)
}

func TestDelete_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	err := service.Delete("test-policy-id")
	assert.Equal(t, errTest, err)
}

func TestDelete_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{UpdateResponse: testdata.OutputDetails}
	service := service.New(repository)

	err := service.Delete("test-policy-id")
	assert.Nil(t, err)
}
