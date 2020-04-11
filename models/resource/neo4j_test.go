package resource_test

import (
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/mocks"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

func TestGetResourcesDatabaseError(t *testing.T) {
	mockSession := mocks.MockSession{
		SetResponse: errorFromDatabase,
	}
	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.Get()
	assert.NotNil(t, err, "Should not be empty")
}

func TestGetResourcesSingleResource(t *testing.T) {
	mockSession := mocks.MockSession{
		SetResponse: validResourceFromDatabase,
	}
	repository := resource.NewNeo4jRepository(mockSession)
	resources, _ := repository.Get()

	var dtos []resource.Element = []resource.Element{}
	dtos = append(dtos, resource.ConstructResourceResponse(resource.Resource{Name: "test-resource", SourceID: "test-source-id"}, "test-id"))
	response := resource.Response{Data: dtos}

	assert.Equal(t, response, resources)
}

func validResourceFromDatabase() (neo4j.Result, error) {
	return mocks.GetMockResult(), nil
}

func errorFromDatabase() (neo4j.Result, error) {
	return nil, errors.New("Database error")
}
