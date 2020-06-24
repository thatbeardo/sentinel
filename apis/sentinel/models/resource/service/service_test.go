package service_test

import (
	"context"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"

	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	policyTestData "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/testdata"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetByIDErrCandidate              string
	GetResponse                      resource.Output
	GetByIDResponse                  resource.OutputDetails
	CreateResponse                   resource.OutputDetails
	UpdateResponse                   resource.OutputDetails
	AssociatePolicyResponse          policy.OutputDetails
	GetAllAssociatedPoliciesResponse policy.Output

	GetErr       error
	GetByIDErr   error
	CreateErr    error
	UpdateErr    error
	AssociateErr error
	DeleteErr    error
}

func (m mockRepository) Get() (resource.Output, error) {
	return m.GetResponse, m.GetErr
}

func (m mockRepository) GetByID(id string) (resource.OutputDetails, error) {
	if m.GetByIDErrCandidate == id {
		return m.GetByIDResponse, models.ErrNotFound
	}
	return m.GetByIDResponse, m.GetByIDErr
}

func (m mockRepository) Create(ctx context.Context, input *resource.Input) (resource.OutputDetails, error) {

	return m.CreateResponse, m.CreateErr
}

func (m mockRepository) AssociatePolicy(string, *policy.Input) (policy.OutputDetails, error) {
	return m.AssociatePolicyResponse, m.AssociateErr
}

func (m mockRepository) Update(resource.Details, *resource.Input) (resource.OutputDetails, error) {
	return m.UpdateResponse, m.UpdateErr
}

func (m mockRepository) GetAllAssociatedPolicies(string) (policy.Output, error) {
	return m.GetAllAssociatedPoliciesResponse, m.AssociateErr
}

func (m mockRepository) Delete(string) error {
	return m.DeleteErr
}

func TestGetServiceDatabaseError(t *testing.T) {
	repository := mockRepository{
		GetErr: models.ErrDatabase,
	}

	service := service.New(repository)
	_, err := service.Get()

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase)
}

func TestGetServiceNoErrors(t *testing.T) {
	repository := mockRepository{
		GetResponse: testdata.Output,
	}

	service := service.New(repository)
	response, err := service.Get()

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, response, testdata.Output, "Response schema does not match")
}

func TestGetByIdServiceNoErrors(t *testing.T) {
	repository := mockRepository{
		GetByIDResponse: testdata.ModificationResult,
	}

	service := service.New(repository)
	response, err := service.GetByID("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, response, testdata.ModificationResult, "Response schema does not match")
}

func TestGetByIdServiceNoErrorsNoResources(t *testing.T) {
	repository := mockRepository{
		GetByIDErr: models.ErrNotFound,
	}

	service := service.New(repository)
	_, err := service.GetByID("test-id")

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Error schemas don't ")
}

func TestGetByIdServiceRepositoryError(t *testing.T) {
	repository := mockRepository{
		GetByIDErr: models.ErrDatabase,
	}

	service := service.New(repository)
	_, err := service.GetByID("test-id")

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Response schema does not match")
}

func TestCreateResourceRepositoryError(t *testing.T) {
	repository := mockRepository{
		GetByIDResponse: testdata.ModificationResult,
		CreateErr:       models.ErrDatabase,
	}

	service := service.New(repository)
	_, err := service.Create(mocks.Context{}, testdata.InputWithoutRelationship)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestCreateResourceNoRelationships(t *testing.T) {
	repository := mockRepository{
		CreateResponse: testdata.ModificationResult,
	}

	service := service.New(repository)
	response, err := service.Create(mocks.Context{}, testdata.InputWithoutRelationship)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, testdata.ModificationResult, response, "Schemas don't match")
}

func TestCreateResourceValidParent(t *testing.T) {
	repository := mockRepository{
		CreateResponse:  testdata.ModificationResult,
		GetByIDResponse: testdata.ModificationResult,
	}

	service := service.New(repository)
	response, err := service.Create(mocks.Context{}, testdata.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, testdata.ModificationResult, response, "Schemas don't match")
}

func TestCreateResourceCreateFails(t *testing.T) {
	repository := mockRepository{
		CreateErr:       models.ErrDatabase,
		GetByIDResponse: testdata.ModificationResult,
	}

	service := service.New(repository)
	_, err := service.Create(mocks.Context{}, testdata.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, models.ErrDatabase, err, "Schemas don't match")
}

func TestDeleteResourceRepositoryError(t *testing.T) {
	repository := mockRepository{
		DeleteErr: models.ErrDatabase,
	}

	service := service.New(repository)
	err := service.Delete("test-id")

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestDeleteResourceRepositoryNoError(t *testing.T) {
	repository := mockRepository{}

	service := service.New(repository)
	err := service.Delete("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
}

func TestUpdateResourceInvalidParent(t *testing.T) {
	repository := mockRepository{
		GetByIDErrCandidate: "parent-id",
		GetByIDResponse:     testdata.ModificationResult,
	}

	service := service.New(repository)
	_, err := service.Update("test-id", testdata.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestUpdateNoErrors(t *testing.T) {
	repository := mockRepository{
		GetByIDResponse: testdata.ModificationResult,
		UpdateResponse:  testdata.ModificationResult,
	}

	service := service.New(repository)
	response, err := service.Update("test-id", testdata.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, testdata.ModificationResult, response, "Response schemas don't match")
}

func TestUpdate_ResourceNodeNotFound_ReturnsErrNotFound(t *testing.T) {
	repository := mockRepository{
		GetByIDErr: models.ErrNotFound,
	}

	service := service.New(repository)
	_, err := service.Update("test-id", testdata.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestAssociatePolicy_RepositoryReturnsError_ReturnError(t *testing.T) {
	repository := mockRepository{
		AssociateErr: models.ErrNotFound,
	}
	service := service.New(repository)
	_, err := service.AssociatePolicy("test-id", policyTestData.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestAssociatePolicy_RepositoryReturnsData_ReturnData(t *testing.T) {
	repository := mockRepository{
		AssociatePolicyResponse: policyTestData.OutputDetails,
	}
	service := service.New(repository)
	results, err := service.AssociatePolicy("test-id", policyTestData.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, policyTestData.OutputDetails, results, "Schemas don't match")
}

func TestGetAllAssociatedPolicies_RepositoryReturnsData_ReturnData(t *testing.T) {
	repository := mockRepository{
		GetAllAssociatedPoliciesResponse: policyTestData.Output,
	}

	service := service.New(repository)
	results, err := service.GetAllAssociatedPolicies("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, policyTestData.Output, results, "Schemas don't match")
}

func TestGetAllAssociatedPolicies_RepositoryReturnsError_ReturnError(t *testing.T) {
	repository := mockRepository{
		AssociateErr: models.ErrDatabase,
	}

	service := service.New(repository)
	_, err := service.GetAllAssociatedPolicies("test-id")

	assert.Equal(t, models.ErrDatabase, err, "Schemas don't match")
}

// func TestUpdateResourceNoParentProvided(t *testing.T) {
// 	repository := &mocks.Repository{}
// 	repository.On("GetByID", m.AnythingOfType("string")).Return(elementWithoutParentNoErrors())
// 	repository.On("GetByID", m.AnythingOfType("string")).Return(parentElementWithoutErrors())
// 	repository.On("Update", m.AnythingOfType("entity.Element"), m.AnythingOfType("*entity.Input")).Return(elementWithoutParentNoErrors())

// 	service := service.NewService(repository)
// 	response, err := service.Update("test-id", data.InputRelationshipsAbsent)

// 	assert.Nil(t, err, "Should not have thrown an error")
// 	assert.Equal(t, data.ElementWithoutParent, response, "Response schemas don't match")
// }
