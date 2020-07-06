package session_test

import (
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/session"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/testdata"
	"github.com/stretchr/testify/assert"
)

type mockNeo4jSession struct {
	RunResponse []map[string]interface{}
	RunErr      error
}

func (m mockNeo4jSession) Run(statement string, parameters map[string]interface{}) ([]map[string]interface{}, error) {
	return m.RunResponse, m.RunErr
}

func TestExecute_RunReturnsError_ReturnDatabaseError(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunErr: errors.New("Database error"),
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_DecodeFails_ReturnDatabaseError(t *testing.T) {
	defer injection.Reset()
	var errDecoding = errors.New("some-decoder-error")

	injection.NodeDecoder = func(map[string]interface{}, string, interface{}) error { return errDecoding }
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: generateValidResultMap(),
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, errDecoding, err)
}

func TestExecute_contextDecodeFails_ReturnDatabaseError(t *testing.T) {
	result := generateValidResultMap()
	result[0]["context"] = "invalid-context"
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: result,
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_NoErrorsFromDB_ReturnResponse(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: generateValidResultMap(),
	})

	response, err := session.Execute(`cypher-query`, map[string]interface{}{})

	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestExecute_DatabaseReturnsNoResources_EmptyResourcesArrayReturned(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: []map[string]interface{}{},
	})

	response, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, testdata.EmptyOutput, response)
	assert.Nil(t, err)
}

func TestExecute_ParentDecodeFails_ReturnError(t *testing.T) {
	results := generateValidResultMap()
	results[0]["parent"] = "invalid entry"

	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: results,
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func generateValidResultMap() []map[string]interface{} {
	result := map[string]interface{}{
		"child": mocks.NewNode(1,
			[]string{"Resource"},
			map[string]interface{}{
				"id":         "test-id",
				"name":       "test-resource",
				"source_id":  "test-source-id",
				"context_id": "test-context-id",
			}),
		"parent": mocks.NewNode(
			2,
			[]string{"Resource"},
			map[string]interface{}{
				"id":        "parent-id",
				"name":      "parent",
				"source_id": "parent-source-id",
			}),
		"context": []interface{}{mocks.NewNode(2,
			[]string{"Resource"},
			map[string]interface{}{
				"id":   "context-id",
				"name": "context",
			}),
		},
	}
	return []map[string]interface{}{result}
}
