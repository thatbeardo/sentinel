package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"github.com/thatbeardo/go-sentinel/mocks"
	"github.com/thatbeardo/go-sentinel/mocks/data"
	models "github.com/thatbeardo/go-sentinel/models"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

func TestGetServiceDatabaseError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Get").Return(errorFromRepository())

	service := resource.NewService(repository)
	_, err := service.Get()

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase)
}

func TestGetServiceNoErrors(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Get").Return(getAllResourcesNoErrors())

	service := resource.NewService(repository)
	response, err := service.Get()

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, response, data.Response, "Response schema does not match")
}

func TestGetByIdServiceNoErrors(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())

	service := resource.NewService(repository)
	response, err := service.GetByID("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, response, data.Element, "Response schema does not match")
}

func TestGetByIdServiceNoErrorsNoResources(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(errorFromRepositoryNotFound())

	service := resource.NewService(repository)
	_, err := service.GetByID("test-id")

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Error schemas don't ")
}

func TestGetByIdServiceRepositoryError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(databaseErrorFromRepository())

	service := resource.NewService(repository)
	_, err := service.GetByID("test-id")

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Response schema does not match")
}

func TestCreateResourceRepositoryError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(databaseErrorFromRepository())

	service := resource.NewService(repository)
	_, err := service.Create(data.InputRelationshipsAbsent)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestCreateResourceParentAbsentInDatabase(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(errorFromRepositoryNotFound())

	service := resource.NewService(repository)
	_, err := service.Create(data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestCreateResourceNoRelationships(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(elementWithoutErrors())

	service := resource.NewService(repository)
	response, err := service.Create(data.InputRelationshipsAbsent)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, data.Element, response, "Schemas don't match")
}

func TestCreateResourceValidParent(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(elementWithoutErrors())
	repository.On("GetByID", "parent-id").Return(parentElementWithoutErrors())
	repository.On("CreateEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(nil)

	service := resource.NewService(repository)
	response, err := service.Create(data.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, data.Element, response, "Schemas don't match")
}

func TestCreateResourceCreateFails(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(databaseErrorFromRepository())
	repository.On("GetByID", "parent-id").Return(parentElementWithoutErrors())

	service := resource.NewService(repository)
	_, err := service.Create(data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, models.ErrDatabase, err, "Schemas don't match")
}

func TestCreateResourceCreateEdgeFailsDatabaseError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(elementWithoutErrors())
	repository.On("GetByID", "parent-id").Return(parentElementWithoutErrors())
	repository.On("CreateEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(models.ErrDatabase)

	service := resource.NewService(repository)
	_, err := service.Create(data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, models.ErrDatabase, err, "Schemas don't match")
}

func TestCreateResourceCreateEdgeFailsNodesNotFound(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(elementWithoutErrors())
	repository.On("GetByID", "parent-id").Return(parentElementWithoutErrors())
	repository.On("CreateEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(models.ErrNotFound)

	service := resource.NewService(repository)
	_, err := service.Create(data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, models.ErrNotFound, err, "Schemas don't match")
}

func TestDeleteResourceRepositoryError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Delete", "test-id").Return(models.ErrDatabase)

	service := resource.NewService(repository)
	err := service.Delete("test-id")

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestDeleteResourceRepositoryNoError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Delete", "test-id").Return(nil)

	service := resource.NewService(repository)
	err := service.Delete("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
}

func TestUpdateResourceInvalidParent(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())
	repository.On("GetByID", "parent-id").Return(errorFromRepositoryNotFound())

	service := resource.NewService(repository)
	_, err := service.Update("test-id", data.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestUpdateResourceValidParentEdgeError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())
	repository.On("CreateEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(models.ErrDatabase)

	service := resource.NewService(repository)
	_, err := service.Update("test-id", data.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestUpdateNoErrors(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())
	repository.On("CreateEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(nil)
	repository.On("DeleteEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(nil)
	repository.On("Update", m.AnythingOfType("string"), m.AnythingOfType("*resource.Input")).Return(elementWithoutErrors())

	service := resource.NewService(repository)
	_, err := service.Update("test-id", data.Input)

	assert.Nil(t, err, "Should not have thrown an error")
}

func TestUpdateResourceNodeNotFound(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(errorFromRepositoryNotFound())

	service := resource.NewService(repository)
	_, err := service.Update("test-id", data.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestUpdateResourceDeleteEdgeError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())
	repository.On("CreateEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(nil)
	repository.On("DeleteEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(models.ErrDatabase)

	service := resource.NewService(repository)
	_, err := service.Update("test-id", data.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestUpdateResourceValidParentDeleteEdgeError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("GetByID", "test-id").Return(elementWithoutRelationshipsNoErrors())
	repository.On("CreateEdge", m.AnythingOfType("string"), m.AnythingOfType("string")).Return(nil)
	repository.On("Update", m.AnythingOfType("string"), m.AnythingOfType("*resource.Input")).Return(elementWithoutErrors())

	service := resource.NewService(repository)
	response, err := service.Update("test-id", data.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, response, data.Element, "Schemas don't match")
}

func getAllResourcesNoErrors() (resource.Response, error) {
	return data.Response, nil
}

func elementWithoutErrors() (resource.Element, error) {
	return data.Element, nil
}

func elementWithoutRelationshipsNoErrors() (resource.Element, error) {
	return data.ElementRelationshipsAbsent, nil
}

func parentElementWithoutErrors() (resource.Element, error) {
	return data.ParentElement, nil
}

func errorFromRepository() (resource.Response, error) {
	return resource.Response{}, models.ErrDatabase
}

func errorFromRepositoryNotFound() (resource.Element, error) {
	return resource.Element{}, models.ErrNotFound
}

func databaseErrorFromRepository() (resource.Element, error) {
	return resource.Element{}, models.ErrDatabase
}
