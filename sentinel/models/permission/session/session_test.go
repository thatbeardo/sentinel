package session_test

import (
	"errors"
	"testing"

	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/session"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/testdata"
	"github.com/stretchr/testify/assert"
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

func TestExecute_DatabaseReturnsNoPolicies_EmptyResourcesArrayReturned(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: []map[string]interface{}{},
	})

	response, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, outputs.Response{Data: []outputs.Permission{}}, response)
	assert.Nil(t, err)
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

	assert.Equal(t, testdata.Response, response)
	assert.Nil(t, err)
}

func generateValidResultMap() []map[string]interface{} {
	result := map[string]interface{}{
		"permission": mockNode{
			id:     1,
			labels: []string{"Permission"},
			props: map[string]interface{}{
				"id":        "test-id",
				"name":      "test-permission",
				"permitted": "allow",
			},
		},
	}
	return []map[string]interface{}{result}
}
