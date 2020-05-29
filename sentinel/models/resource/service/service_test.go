package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"github.com/thatbeardo/go-sentinel/mocks"
	"github.com/thatbeardo/go-sentinel/mocks/data"
	models "github.com/thatbeardo/go-sentinel/models"
	entity "github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/models/resource/service"
)

func TestGetServiceDatabaseError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Get").Return(errorFromRepository())

	service := service.NewService(repository)
	_, err := service.Get()

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase)
}

func TestGetServiceNoErrors(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Get").Return(getAllResourcesNoErrors())

	service := service.NewService(repository)
	response, err := service.Get()

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, response, data.Response, "Response schema does not match")
}

func TestGetByIdServiceNoErrors(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())

	service := service.NewService(repository)
	response, err := service.GetByID("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, response, data.Element, "Response schema does not match")
}

func TestGetByIdServiceNoErrorsNoResources(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(errorFromRepositoryNotFound())

	service := service.NewService(repository)
	_, err := service.GetByID("test-id")

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Error schemas don't ")
}

func TestGetByIdServiceRepositoryError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(databaseErrorFromRepository())

	service := service.NewService(repository)
	_, err := service.GetByID("test-id")

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Response schema does not match")
}

func TestCreateResourceRepositoryError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("Create", m.AnythingOfType("*entity.Input")).Return(databaseErrorFromRepository())

	service := service.NewService(repository)
	_, err := service.Create(data.InputRelationshipsAbsent)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestCreateResourceParentAbsentInDatabase(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(errorFromRepositoryNotFound())

	service := service.NewService(repository)
	_, err := service.Create(data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestCreateResourceNoRelationships(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*entity.Input")).Return(elementWithoutErrors())

	service := service.NewService(repository)
	response, err := service.Create(data.InputRelationshipsAbsent)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, data.Element, response, "Schemas don't match")
}

func TestCreateResourceValidParent(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*entity.Input")).Return(elementWithoutErrors())
	repository.On("GetByID", "parent-id").Return(parentElementWithoutErrors())

	service := service.NewService(repository)
	response, err := service.Create(data.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, data.Element, response, "Schemas don't match")
}

func TestCreateResourceCreateFails(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Create", m.AnythingOfType("*entity.Input")).Return(databaseErrorFromRepository())
	repository.On("GetByID", "parent-id").Return(parentElementWithoutErrors())

	service := service.NewService(repository)
	_, err := service.Create(data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, models.ErrDatabase, err, "Schemas don't match")
}

func TestDeleteResourceRepositoryError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Delete", "test-id").Return(models.ErrDatabase)

	service := service.NewService(repository)
	err := service.Delete("test-id")

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestDeleteResourceRepositoryNoError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("Delete", "test-id").Return(nil)

	service := service.NewService(repository)
	err := service.Delete("test-id")

	assert.Nil(t, err, "Should not have thrown an error")
}

func TestUpdateResourceInvalidParent(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())
	repository.On("GetByID", "parent-id").Return(errorFromRepositoryNotFound())

	service := service.NewService(repository)
	_, err := service.Update("test-id", data.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestUpdateNoErrors(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("GetByID", "test-id").Return(elementWithoutErrors())
	repository.On("Update", m.AnythingOfType("entity.Element"), m.AnythingOfType("*entity.Input")).Return(elementWithoutErrors())

	service := service.NewService(repository)
	response, err := service.Update("test-id", data.Input)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, data.Element, response, "Response schemas don't match")
}

func TestUpdateResourceNodeNotFound(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "test-id").Return(errorFromRepositoryNotFound())

	service := service.NewService(repository)
	_, err := service.Update("test-id", data.Input)

	assert.NotNil(t, err, "Should not have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
}

func TestUpdateResourceNoParentProvided(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", m.AnythingOfType("string")).Return(elementWithoutParentNoErrors())
	repository.On("GetByID", m.AnythingOfType("string")).Return(parentElementWithoutErrors())
	repository.On("Update", m.AnythingOfType("entity.Element"), m.AnythingOfType("*entity.Input")).Return(elementWithoutParentNoErrors())

	service := service.NewService(repository)
	response, err := service.Update("test-id", data.InputRelationshipsAbsent)

	assert.Nil(t, err, "Should not have thrown an error")
	assert.Equal(t, data.ElementWithoutParent, response, "Response schemas don't match")
}

func getAllResourcesNoErrors() (entity.Response, error) {
	return data.Response, nil
}

func elementWithoutErrors() (entity.Element, error) {
	return data.Element, nil
}

func elementWithoutRelationshipsNoErrors() (entity.Element, error) {
	return data.ElementRelationshipsAbsent, nil
}

func elementWithoutParentNoErrors() (entity.Element, error) {
	return data.ElementWithoutParent, nil
}

func parentElementWithoutErrors() (entity.Element, error) {
	return data.ParentElement, nil
}

func errorFromRepository() (entity.Response, error) {
	return entity.Response{}, models.ErrDatabase
}

func errorFromRepositoryNotFound() (entity.Element, error) {
	return entity.Element{}, models.ErrNotFound
}

func databaseErrorFromRepository() (entity.Element, error) {
	return entity.Element{}, models.ErrDatabase
}
