package resource_test

import (
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/mocks"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

func TestGetResourcesOk(t *testing.T) {
	mockSession := mocks.MockSession{
		SetResponse: errorFromDatabase,
	}
	repository := resource.NewNeo4jRepository(mockSession)
	_, err := repository.Get()
	assert.NotNil(t, err, "Should not be empty")
}

func errorFromDatabase() (neo4j.Result, error) {
	return nil, errors.New("Database error")
}
