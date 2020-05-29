package session_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/bithippie/guard-my-app/sentinel/mocks/data"
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/session"
)

type mockNode struct {
	id     int64
	labels []string
	props  map[string]interface{}
}

func (m mockNode) Id() int64 {
	return m.id
}

func (m mockNode) Labels() []string {
	return m.labels
}

func (m mockNode) Props() map[string]interface{} {
	return m.props
}

type mockNeo4jSession struct {
	RunResponse map[string]interface{}
	RunErr      error
}

func (m mockNeo4jSession) Run(statement string, parameters map[string]interface{}) (map[string]interface{}, error) {
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

	injection.MapDecoder = func(interface{}, interface{}) error { return errDecoding }
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: generateValidResultMap(),
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, errDecoding, err)
}

func TestExecute_NoErrorsFromDB_ReturnResponse(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: generateValidResultMap(),
	})

	response, err := session.Execute(`cypher-query`, map[string]interface{}{})

	assert.Equal(t, data.ResponseWithoutPolicies, response)
	assert.Nil(t, err)
}

func TestExecute_ParentDecodeFails_ReturnError(t *testing.T) {
	resultMap := generateValidResultMap()
	resultMap["parent"] = "invalid entry"

	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: resultMap,
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func generateValidResultMap() map[string]interface{} {
	return map[string]interface{}{
		"child": mockNode{
			id:     1,
			labels: []string{"Resource"},
			props: map[string]interface{}{
				"id":        "test-id",
				"name":      "test-resource",
				"source_id": "test-source-id",
			},
		},
		"parent": mockNode{
			id:     2,
			labels: []string{"Resource"},
			props: map[string]interface{}{
				"id":        "parent-id",
				"name":      "parent",
				"source_id": "parent-source-id",
			},
		},
	}
}
