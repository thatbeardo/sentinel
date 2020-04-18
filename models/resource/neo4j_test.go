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

	var dtos []resource.Element = []resource.Element{data.ElementWithoutParent}
	response := resource.Response{Data: dtos}

	assert.Nil(t, err, "Should be empty")
	assert.Equal(t, response, resources, "Response schemas don't match")
}

func TestGetResourcesResourceWithParent(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(validResourceFromDatabaseWithParent())

	repository := resource.NewNeo4jRepository(mockSession)
	resources, err := repository.Get()

	var dtos []resource.Element = []resource.Element{data.Element}
	response := resource.Response{Data: dtos}

	assert.Nil(t, err, "Should be empty")
	assert.Equal(t, response, resources, "Response schemas don't match")
}

func TestGetResourcesByIdSingleResource(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(validResourceFromDatabaseWithParentGetByID())

	repository := resource.NewNeo4jRepository(mockSession)
	resources, err := repository.GetByID("test-id")
	response := data.ElementWithoutParent

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
	element, err := repository.Create(data.Input)

	assert.Nil(t, err, "Error should be empty")
	assert.Equal(t, data.ElementWithoutParent, element, "Error model does not match")
}

func TestCreateResourcesDatabaseError(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(errorFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.Create(data.Input)

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

func TestDeleteResourcesSummaryFailure(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(summaryFailure())

	repository := resource.NewNeo4jRepository(mockSession)
	err := repository.Delete("test-id")

	assert.Equal(t, err, models.ErrDatabase, "Error schemas do not match")
}

func TestDeleteResourcesNoNodesDeleted(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(deleteResourceNoNodesDeleted())

	repository := resource.NewNeo4jRepository(mockSession)
	err := repository.Delete("test-id")

	assert.Equal(t, err, models.ErrNotFound, "Error schemas do not match")
}

func TestUpdateResourceDatabaseError(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(errorFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.Update("test-id", data.Input)

	assert.NotNil(t, err)
	assert.Equal(t, err, models.ErrDatabase, "Error schemas do not match")
}

func TestUpdateResourceNoErrors(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(createResourcesSuccessful())

	repository := resource.NewNeo4jRepository(mockSession)
	response, err := repository.Update("test-id", data.Input)

	assert.Nil(t, err)
	assert.Equal(t, response, data.Element, "Error schemas do not match")
}

func TestUpdateOwnershipNoErrors(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(errorFromDatabase())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.UpdateOwnership("test-id", data.Input)

	assert.NotNil(t, err)
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestUpdateOwnershipSummaryErrors(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(summaryFailure())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.UpdateOwnership("test-id", data.Input)

	assert.NotNil(t, err)
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestUpdateOwnershipZeroRelationshipsCreated(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(updateOwnershipZeroRelationsCreated())

	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.UpdateOwnership("test-id", data.Input)

	assert.NotNil(t, err)
	assert.Equal(t, err, models.ErrDatabase, "Schemas don't match")
}

func TestUpdateOwnershipNoErros(t *testing.T) {
	mockSession := &mocks.Session{}
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(UpdateOwnershipNoErrors())

	repository := resource.NewNeo4jRepository(mockSession)
	response, err := repository.UpdateOwnership("test-id", data.Input)

	assert.Nil(t, err)
	assert.Equal(t, response, data.Element, "Schemas don't match")
}

func createResourcesSuccessful() (neo4j.Result, error) {
	return mocks.CreateResourceSuccessful(), nil
}

func noResourcesFromDatabase() (neo4j.Result, error) {
	return mocks.GetEmptyResult(), nil
}

func validResourceFromDatabase() (neo4j.Result, error) {
	return mocks.ResourceWithoutParent(), nil
}

func validResourceFromDatabaseWithParent() (neo4j.Result, error) {
	return mocks.ResourceWithParent(), nil
}

func validResourceFromDatabaseWithParentGetByID() (neo4j.Result, error) {
	return mocks.ResourceWithoutParentGetByID(), nil
}

func deleteResourceSuccessful() (neo4j.Result, error) {
	return mocks.DeleteResourceSuccessful(), nil
}

func deleteRelationshipSuccessful() (neo4j.Result, error) {
	return mocks.DeleteRelationshipSuccessful(), nil
}

func deleteResourceNoNodesDeleted() (neo4j.Result, error) {
	return mocks.DeleteResourceNoNodesDeleted(), nil
}

func summaryFailure() (neo4j.Result, error) {
	return mocks.SummaryFailure(), nil
}

func errorFromDatabase() (neo4j.Result, error) {
	return nil, errors.New("Database error")
}

func updateOwnershipZeroRelationsCreated() (neo4j.Result, error) {
	return mocks.UpdateOwnershipZeroRelationshipsCreated(), nil
}

func UpdateOwnershipNoErrors() (neo4j.Result, error) {
	return mocks.UpdateOwnershipNoErrors(), nil
}
