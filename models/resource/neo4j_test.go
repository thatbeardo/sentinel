package resource_test

import (
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thatbeardo/go-sentinel/mocks"
	"github.com/thatbeardo/go-sentinel/mocks/data"
	models "github.com/thatbeardo/go-sentinel/models"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

func TestGetResourcesDatabaseError(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(errorFromDatabase())
	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.Get()
	assert.Equal(t, models.ErrDatabase, err, "Error model does not match")
	assert.NotNil(t, err, "Should not be empty")
}

func TestGetResourcesNoResourcesPresent(t *testing.T) {
	mockSession := &mocks.Session{}
	mockResult := &mocks.Result{}

	mockResult.On("Next").Return(false).Once()
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(noResourcesFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	resources, err := repository.Get()

	var dtos []resource.Element = []resource.Element{}
	response := resource.Response{Data: dtos}

	assert.Nil(t, err, "Should be empty")
	assert.Equal(t, response, resources, "Schemas do not match")
}

func TestGetResourcesSingleResource(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(validResourceFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	resources, err := repository.Get()

	var dtos []resource.Element = []resource.Element{data.Element}
	response := resource.Response{Data: dtos}

	assert.Nil(t, err, "Should be empty")
	assert.Equal(t, response, resources, "Response schemas don't match")
}

func TestGetResourcesByIdSingleResource(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(validResourceFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	resources, err := repository.GetByID("test-id")

	response := data.Element

	assert.Nil(t, err, "Should be empty")
	assert.Equal(t, response, resources, "Response schemas don't match")
}

func TestGetResourcesByIdDatabaseError(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(errorFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.GetByID("test-id")

	assert.NotNil(t, err, "Should not be empty")
	assert.Equal(t, models.ErrDatabase, err, "Error model does not match")
}

func TestGetResourcesByIdResourceNotFound(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(noResourcesFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.GetByID("test-id")

	assert.NotNil(t, err, "Should not be empty")
	assert.Equal(t, models.ErrNotFound, err, "Error model does not match")
}

func TestCreateResourcesSuccessful(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(createResourcesSuccessful())

	repository := resource.NewNeo4jRepository(mockSession)
	element, err := repository.Create(&data.Input)

	assert.Nil(t, err, "Error should be empty")
	assert.Equal(t, data.Element, element, "Error model does not match")
}

func TestCreateResourcesDatabaseError(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(errorFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.Create(&data.Input)

	assert.NotNil(t, err, "Error should be empty")
	assert.Equal(t, models.ErrDatabase, err, "Error model does not match")
}

func TestDeleteResourcesDatabaseError(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(errorFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	err := repository.Delete("test-id")

	assert.NotNil(t, err, "Error should be empty")
	assert.Equal(t, models.ErrDatabase, err, "Error model does not match")
}

func TestDeleteResourceSuccessful(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(deleteResourceSuccessful())

	repository := resource.NewNeo4jRepository(mockSession)
	err := repository.Delete("test-id")

	assert.Nil(t, err, "Error should be empty")
}

func TestDeleteResourcesNoNodesDeleted(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(deleteResourceNoNodesDeleted())

	repository := resource.NewNeo4jRepository(mockSession)
	err := repository.Delete("test-id")

	assert.Equal(t, err, models.ErrNotFound, "Error schemas do not match")
}

// func TestCreateEdgeNoErrors(t *testing.T) {
// 	mockSession := &mocks.Session{}
// 	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return()

// 	repository := resource.NewNeo4jRepository(mockSession)
// 	err := repository.Delete("test-id")

// 	assert.Equal(t, err, models.ErrNotFound, "Error schemas do not match")
// }

func createResourcesSuccessful() (neo4j.Result, error) {
	return mocks.CreateResourceSuccessful(), nil
}

func noResourcesFromDatabase() (neo4j.Result, error) {
	return mocks.GetEmptyResult(), nil
}

func validResourceFromDatabase() (neo4j.Result, error) {
	return mocks.GetMockResult(), nil
}

func deleteResourceSuccessful() (neo4j.Result, error) {
	return mocks.DeleteResourceSuccessful(), nil
}

func deleteResourceNoNodesDeleted() (neo4j.Result, error) {
	return mocks.DeleteResourceNoNodesDeleted(), nil
}

func errorFromDatabase() (neo4j.Result, error) {
	return nil, errors.New("Database error")
}
