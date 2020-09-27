package service_test

import (
	"errors"
	"testing"

	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetAllPermissionsForcontextResponse             permission.Output
	GetAllPermissionsForcontextResponseWithResource permission.Output
	CreateResponse                                  permission.OutputDetails
	UpdateResponse                                  permission.OutputDetails
	Err                                             error
}

func (m mockRepository) Create(*permission.Input, string, string) (permission.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockRepository) GetAllPermissionsForcontext(contextID string) (permission.Output, error) {
	return m.GetAllPermissionsForcontextResponse, m.Err
}

func (m mockRepository) GetAllPermissionsForcontextWithResource(contextID string, resourceID string) (permission.Output, error) {
	return m.GetAllPermissionsForcontextResponseWithResource, m.Err
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

	_, err := service.Create(testdata.Input, "test-context-id", "test-target-id")
	assert.Equal(t, errTest, err)
}

func TestCreate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{CreateResponse: testdata.OutputDetails}
	service := service.New(repository)

	permission, err := service.Create(testdata.Input, "test-context-id", "test-target-id")
	assert.Equal(t, testdata.OutputDetails, permission)
	assert.Nil(t, err)
}

func TestGetAllPermissions_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.GetAllPermissionsForcontext("test-context-id")
	assert.Equal(t, errTest, err)
}

func TestGetAllPermissions_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{GetAllPermissionsForcontextResponse: testdata.Output}
	service := service.New(repository)

	permission, err := service.GetAllPermissionsForcontext("test-context-id")
	assert.Equal(t, testdata.Output, permission)
	assert.Nil(t, err)
}

func TestGetAllPermissionsOfResource_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.GetAllPermissionsForcontextWithResource("test-context-id", "test-resource-id")
	assert.Equal(t, errTest, err)
}

func TestGetAllPermissionsOfResource_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{GetAllPermissionsForcontextResponseWithResource: testdata.Output}
	service := service.New(repository)

	permission, err := service.GetAllPermissionsForcontextWithResource("test-context-id", "test-resource-id")
	assert.Equal(t, testdata.Output, permission)
	assert.Nil(t, err)
}

func TestUpdate_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	_, err := service.Update("test-context-id", testdata.Input)
	assert.Equal(t, errTest, err)
}

func TestUpdate_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{UpdateResponse: testdata.OutputDetails}
	service := service.New(repository)

	permission, err := service.Update("test-context-id", testdata.Input)
	assert.Equal(t, testdata.OutputDetails, permission)
	assert.Nil(t, err)
}

func TestDelete_RepositoryReturnsError_ReturnsError(t *testing.T) {
	repository := mockRepository{Err: errTest}
	service := service.New(repository)

	err := service.Delete("test-context-id")
	assert.Equal(t, errTest, err)
}

func TestDelete_RepositoryReturnsResponse_ResponseReturned(t *testing.T) {
	repository := mockRepository{UpdateResponse: testdata.OutputDetails}
	service := service.New(repository)

	err := service.Delete("test-context-id")
	assert.Nil(t, err)
}

func TestPermissionIdempotence_DuplicatePermissionName_ReturnsFalse(t *testing.T) {
	repository := mockRepository{
		GetAllPermissionsForcontextResponseWithResource: testdata.Output,
	}
	service := service.New(repository)
	isPermissionIdempotent, _ := service.IsPermissionIdempotent(testdata.Input, "test-context-id", "test-principal-id")
	assert.False(t, isPermissionIdempotent)
}

func TestPermissionIdempotence_UniquePermissionName_ReturnsTrue(t *testing.T) {
	repository := mockRepository{
		GetAllPermissionsForcontextResponseWithResource: testdata.Output,
	}
	service := service.New(repository)

	idempotentPermission := clonePermission(testdata.Input)
	idempotentPermission.Data.Attributes.Name = "permission"

	isPermissionIdempotent, _ := service.IsPermissionIdempotent(idempotentPermission, "test-context-id", "test-principal-id")
	assert.True(t, isPermissionIdempotent)
}

func TestPermissionIdempotence_SamePermissionNameDifferentPermitted_ReturnsFalse(t *testing.T) {
	repository := mockRepository{
		GetAllPermissionsForcontextResponseWithResource: testdata.Output,
	}
	service := service.New(repository)

	idempotentPermission := clonePermission(testdata.Input)
	idempotentPermission.Data.Attributes.Permitted = "deny"

	isPermissionIdempotent, _ := service.IsPermissionIdempotent(idempotentPermission, "test-context-id", "test-principal-id")
	assert.False(t, isPermissionIdempotent)
}

func clonePermission(original *permission.Input) (input *permission.Input) {
	input = &permission.Input{
		Data: permission.InputDetails{
			Type: "permission",
			Attributes: &permission.Attributes{
				Name:      original.Data.Attributes.Name,
				Permitted: original.Data.Attributes.Permitted,
			},
		},
	}
	return
}
