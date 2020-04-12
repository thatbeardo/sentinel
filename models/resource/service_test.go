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

func TestCreateResourceNoErrors(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(elementWithoutErrors())

	service := resource.NewService(repository)
	response, err := service.Create(&data.Input)

	assert.Nil(t, err, "Shouldn't have thrown any errors")
	assert.Equal(t, response, data.Element, "Schemas don't match")
}

func TestCreateResourceRepositoryError(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(elementWithoutErrors())
	repository.On("Create", m.AnythingOfType("*resource.Input")).Return(databaseErrorFromRepository())

	service := resource.NewService(repository)
	_, err := service.Create(&data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestCreateResourceParentAbsent(t *testing.T) {
	repository := &mocks.Repository{}
	repository.On("GetByID", "parent-id").Return(errorFromRepositoryNotFound())

	service := resource.NewService(repository)
	_, err := service.Create(&data.Input)

	assert.NotNil(t, err, "Should have thrown an error")
	assert.Equal(t, err, models.ErrNotFound, "Schemas don't match")
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

func getAllResourcesNoErrors() (resource.Response, error) {
	return data.Response, nil
}

func elementWithoutErrors() (resource.Element, error) {
	return data.Element, nil
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
