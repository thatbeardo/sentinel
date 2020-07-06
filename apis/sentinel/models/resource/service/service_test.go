package service_test

import (
	"context"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"

	contextDto "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	contextTestData "github.com/bithippie/guard-my-app/apis/sentinel/models/context/testdata"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/testdata"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	GetByIDErrCandidate              string
	GetResponse                      resource.Output
	GetByIDResponse                  resource.OutputDetails
	GetResourcesHubResponse          resource.OutputDetails
	CreateTenantResponse             resource.OutputDetails
	CreateResourceResponse           resource.OutputDetails
	UpdateResponse                   resource.OutputDetails
	UpdateDefaultResponse            resource.OutputDetails
	AssociatecontextResponse         contextDto.OutputDetails
	GetAllAssociatedContextsResponse contextDto.Output

	GetErr             error
	GetByIDErr         error
	GetResourcesHubErr error
	CreateErr          error
	UpdateErr          error
	AssociateErr       error
	DeleteErr          error

	T *testing.T
}

func (m mockRepository) Get(string, string) (resource.Output, error) {
	return m.GetResponse, m.GetErr
}

func (m mockRepository) GetByID(id string) (resource.OutputDetails, error) {
	if m.GetByIDErrCandidate == id {
		return m.GetByIDResponse, models.ErrNotFound
	}
	return m.GetByIDResponse, m.GetByIDErr
}

func (m mockRepository) GetResourcesHub(clientID, tenant string) (resource.OutputDetails, error) {
	return m.GetResourcesHubResponse, m.GetResourcesHubErr
}

func (m mockRepository) Create(input *resource.Input) (resource.OutputDetails, error) {
	assert.Equal(m.T, testdata.Input, input)
	return m.CreateResourceResponse, m.CreateErr
}

func (m mockRepository) CreateMetaResource(*resource.Input) (resource.OutputDetails, error) {
	return m.CreateTenantResponse, m.CreateErr
}

func (m mockRepository) AddContext(string, *contextDto.Input) (contextDto.OutputDetails, error) {
	return m.AssociatecontextResponse, m.AssociateErr
}

func (m mockRepository) Update(resource.Details, *resource.Input) (resource.OutputDetails, error) {
	return m.UpdateResponse, m.UpdateErr
}

func (m mockRepository) UpdateDefaultContext(string, string) (resource.OutputDetails, error) {
	return m.UpdateDefaultResponse, m.UpdateErr
}

func (m mockRepository) GetAllContexts(string) (contextDto.Output, error) {
	return m.GetAllAssociatedContextsResponse, m.AssociateErr
}

func (m mockRepository) Delete(string) error {
	return m.DeleteErr
}

func TestGetServiceDatabaseError(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaims = func(ctx context.Context, claim string) string {
		return "test-tenant"
	}
	injection.ExtractTenant = func(ctx context.Context) string {
		return "tenant"
	}
	repository := mockRepository{
		GetErr: models.ErrDatabase,
	}

	service := service.New(repository)
	_, err := service.Get(mocks.Context{})

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase)
}

func TestGetServiceNoErrors(t *testing.T) {
	injection.ExtractClaims = func(ctx context.Context, claim string) string {
		return "test-tenant"
	}
	repository := mockRepository{
		GetResponse: testdata.Output,
	}

	service := service.New(repository)
	response, err := service.Get(mocks.Context{})

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

func TestCreateTenantResource_RepositoryReturnsData_ParseAndReturnData(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaims = func(ctx context.Context, claim string) string {
		return "create:resource"
	}
	repository := mockRepository{
		CreateTenantResponse: testdata.ModificationResult,
		T:                    t,
	}
	service := service.New(repository)
	response, err := service.Create(mocks.Context{}, testdata.InputWithoutRelationship)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, testdata.ModificationResult, response)
}

func TestCreate_RelationshipsAbsent_ShouldCallRepositoryWithResourceHub(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaims = func(ctx context.Context, claim string) string {
		return ""
	}

	injection.ExtractTenant = func(ctx context.Context) string {
		return "dev"
	}
	repository := mockRepository{
		GetResourcesHubResponse: testdata.ResourcesHubOutputDetails,
		CreateResourceResponse:  testdata.OutputDetails,
		T:                       t,
	}

	service := service.New(repository)
	response, err := service.Create(mocks.Context{}, testdata.InputWithoutRelationship)

	assert.Nil(t, err, "Should have thrown an error")
	assert.Equal(t, testdata.OutputDetails, response, "Schemas don't match")
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

func TestAssociatecontext_RepositoryReturnsError_ReturnError(t *testing.T) {
	repository := mockRepository{
		AssociateErr: models.ErrNotFound,
	}
	service := service.New(repository)
	_, err := service.Associatecontext("test-id", contextTestData.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestAssociatecontext_RepositoryReturnsData_ReturnData(t *testing.T) {
	repository := mockRepository{
		AssociatecontextResponse: contextTestData.OutputDetails,
	}
	service := service.New(repository)
	results, err := service.Associatecontext("test-id", contextTestData.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, contextTestData.OutputDetails, results, "Schemas don't match")
}

func TestGetAllAssociatedContexts_RepositoryReturnsData_ReturnData(t *testing.T) {
	repository := mockRepository{
		GetAllAssociatedContextsResponse: contextTestData.Output,
	}

	service := service.New(repository)
	results, err := service.GetAllAssociatedContexts("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, contextTestData.Output, results, "Schemas don't match")
}

func TestGetAllAssociatedContexts_RepositoryReturnsError_ReturnError(t *testing.T) {
	repository := mockRepository{
		AssociateErr: models.ErrDatabase,
	}

	service := service.New(repository)
	_, err := service.GetAllAssociatedContexts("test-id")

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
